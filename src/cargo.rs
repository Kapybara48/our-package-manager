use serde::Deserialize;

use crate::error::OurError;

#[derive(Deserialize)]
struct Config {
    package: Package,
}

#[derive(Deserialize)]
struct Package {
    name: String,
}

pub fn get_package_name(package_dir: &std::path::PathBuf) -> Result<String, OurError> {
    let mut config_path = package_dir.clone();
    config_path.push("Cargo.toml");
    let config_file = std::fs::read_to_string(config_path)?;

    let config: Config = toml::from_str(config_file.as_str())?;

    Ok(config.package.name)
}
