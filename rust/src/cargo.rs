use serde::Deserialize;
use std::path::{Path, PathBuf};

use crate::error::OurError;
use crate::{package_config, paths};

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

pub fn get_source_binary_path(package_dir: &Path) -> Result<PathBuf, OurError> {
    let package_name = get_package_name(package_dir)?;
    let mut bin_path = package_dir.to_path_buf();
    bin_path.push("target");
    bin_path.push("debug");
    bin_path.push(package_name);

    Ok(bin_path)
}

pub fn get_destination_binary_path(package_dir: &Path) -> Result<PathBuf, OurError> {
    let mut destination_bin_path = paths::get_bin_dir()?;
    let package_name = get_package_name(package_dir)?;
    destination_bin_path.push(package_name);

    Ok(destination_bin_path)
}

pub fn install_binary(config: &package_config::Config) -> Result<PathBuf, OurError> {
    let source_bin_path = &config.install.binary_source;
    let destination_bin_path = &config.install.binary_destination;

    std::fs::copy(source_bin_path, destination_bin_path)?;

    Ok(destination_bin_path.to_path_buf())
}
