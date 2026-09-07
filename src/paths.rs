use crate::error::OurError;
use std::path::PathBuf;

pub fn get_our_dir() -> Result<PathBuf, OurError> {
    match std::env::home_dir() {
        Some(mut path) => {
            path.push(".our");
            Ok(path)
        }
        None => Err(OurError::MissingHomeEnv),
    }
}
