use serde::Deserialize;

use crate::error::OurError;
use std::path::Path;

#[derive(Deserialize)]
pub struct Config {
    package: Package,
    version: Version,
    source: Source,
    git: Git,
    build: Build,
    install: Install,
}

#[derive(Deserialize)]
pub struct Package {
    name: String,
    path: String,
}

#[derive(Deserialize)]
pub struct Version {
    commit: String,
}

#[derive(Deserialize)]
pub struct Source {
    url: String,
}

#[derive(Deserialize)]
pub struct Git {
    clone_depth: u32,
    branch: String,
}

#[derive(Deserialize)]
pub struct Build {
    command: String,
    args: Vec<String>,
}

#[derive(Deserialize)]
struct Install {
    binary_source: String,
    binary_destination: String,
}

pub fn load_config(package_dir: &Path) -> Result<Option<Config>, OurError> {
    let mut our_path = package_dir.to_path_buf();
    our_path.push("our.toml");

    if !our_path.exists() {
        return Ok(None);
    }

    let config_file = std::fs::read_to_string(our_path)?;

    let config: Config = toml::from_str(&config_file)?;

    Ok(Some(config))
}
