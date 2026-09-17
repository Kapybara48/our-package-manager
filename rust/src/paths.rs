use crate::error::OurError;
use std::fs::create_dir_all;
use std::path::PathBuf;

pub fn get_our_dir() -> Result<PathBuf, OurError> {
    match std::env::home_dir() {
        Some(mut path) => {
            path.push(".our");
            create_dir_all(&path)?;
            Ok(path)
        }
        None => Err(OurError::MissingHomeEnv),
    }
}

pub fn get_bin_dir() -> Result<PathBuf, OurError> {
    let mut our_bin_dir = get_our_dir()?;
    our_bin_dir.push("bin");
    create_dir_all(&our_bin_dir)?;

    Ok(our_bin_dir)
}

pub fn get_packages_dir() -> Result<PathBuf, OurError> {
    let mut our_package_dir = get_our_dir()?;
    our_package_dir.push("packages");
    create_dir_all(&our_package_dir)?;

    Ok(our_package_dir)
}

pub fn get_temp_dir() -> Result<PathBuf, OurError> {
    let mut our_temp_dir = get_our_dir()?;
    our_temp_dir.push("temp");
    create_dir_all(&our_temp_dir)?;

    Ok(our_temp_dir)
}
