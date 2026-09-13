use serde::Deserialize;
use std::path::{Path, PathBuf};

use crate::error::OurError;

#[derive(Deserialize)]
struct Config {
    package: Package,
}

#[derive(Deserialize)]
struct Package {
    name: String,
}

pub fn get_package_name(package_dir: &Path) -> Result<String, OurError> {
    let mut config_path = package_dir.to_path_buf();
    config_path.push("Cargo.toml");
    let config_file = std::fs::read_to_string(config_path)?;

    let config: Config = toml::from_str(config_file.as_str())?;

    Ok(config.package.name)
}

pub fn get_binary_path(package_dir: &Path) -> Result<PathBuf, OurError> {
    let package_name = get_package_name(package_dir)?;
    let mut binary_path = package_dir.to_path_buf();
    binary_path.push("target");
    binary_path.push("debug");
    binary_path.push(package_name);

    Ok(binary_path)
}
