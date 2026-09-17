use crate::{error::OurError, paths::get_our_dir};
use std::{path::PathBuf, process::Command};

pub fn clone(url: &str) -> Result<PathBuf, OurError> {
    let package_name = get_name_from_url(url);
    let mut package_dir = get_our_dir()?;

    let status = Command::new("git")
        .args(["clone", url, package_name])
        .current_dir(&package_dir)
        .spawn()?
        .wait()?;

    if !status.success() {
        return Err(OurError::CloneFailed);
    }

    package_dir.push(package_name);

    Ok(package_dir)
}

fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
