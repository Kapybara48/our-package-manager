use serde::Deserialize;

use crate::{cargo, error::OurError, project::ProjectType};
use std::path::{Path, PathBuf};

#[derive(Deserialize)]
pub struct Config {
    package: Package,
    version: Option<Version>,
    source: Source,
    git: Git,
    pub build: Build,
    pub install: Install,
}

#[derive(Deserialize)]
pub struct Package {
    name: String,
    path: Option<String>,
}

#[derive(Deserialize)]
pub struct Version {
    commit: Option<String>,
}

#[derive(Deserialize)]
pub struct Source {
    url: String,
}

#[derive(Deserialize)]
pub struct Git {
    clone_depth: Option<u32>,
    branch: Option<String>,
}

#[derive(Deserialize)]
pub struct Build {
    pub command: String,
    pub args: Vec<String>,
}

#[derive(Deserialize)]
pub struct Install {
    pub binary_source: PathBuf,
    pub binary_destination: PathBuf,
}

pub struct PackageInfo {
    pub name: String,
    pub url: String,
    pub path: Option<String>,
    pub branch: Option<String>,
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

pub fn generate_config(
    project_type: ProjectType,
    package_info: PackageInfo,
    package_dir: &Path,
) -> Result<Config, OurError> {
    match project_type {
        ProjectType::Rust => Ok(Config {
            package: Package {
                name: package_info.name,
                path: package_info.path,
            },
            version: None,
            source: Source {
                url: package_info.url,
            },
            git: Git {
                branch: package_info.branch,
                clone_depth: Some(1),
            },
            build: Build {
                command: "cargo".to_string(),
                args: vec!["build".to_string()],
            },
            install: Install {
                binary_source: cargo::get_source_binary_path(package_dir)?,
                binary_destination: cargo::get_destination_binary_path(package_dir)?,
            },
        }),
        ProjectType::Unknown => Err(OurError::UnknownProjectType),
    }
}
