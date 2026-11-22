use crate::error::OurError;
use std::fs::create_dir_all;
use std::path::{PathBuf, Path};

pub fn get_home_dir() -> Result<PathBuf, OurError>{
    match std::env::home_dir() {
        Some(home) => Ok(home),
        None => Err(OurError::MissingHomeEnv),
    }
}

pub fn get_our_dir() -> Result<PathBuf, OurError> {
    let mut our_dir = get_home_dir()?;
    our_dir.push(".our");
    create_dir_all(&our_dir)?;

    Ok(our_dir)
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

pub fn clear_temp_dir() -> Result<(), OurError> {
    let temp_dir = get_temp_dir()?;

    for entry in std::fs::read_dir(&temp_dir)? {
        let entry = entry?;
        std::fs::remove_dir_all(entry.path())?;
    }

    Ok(())
}

pub fn resolve_path(path: &Path) -> Result<PathBuf, OurError>{
    if path.starts_with("~"){
        let mut home = get_home_dir()?;
        let relative = path.strip_prefix("~").unwrap();
        home.push(relative);

        Ok(home)
    }else{
        Ok(path.to_path_buf())
    }
}
