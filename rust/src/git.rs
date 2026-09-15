use crate::{error::OurError, paths::get_our_dir};
use std::{
    path::PathBuf,
    process::{Command, ExitStatus},
};

pub fn clone(url: &str) -> Result<(ExitStatus, PathBuf), OurError> {
    let package_name = get_name_from_url(url);
    let mut package_dir = get_our_dir()?;

    let status = Command::new("git")
        .args(["clone", url, package_name])
        .current_dir(&package_dir)
        .spawn()?
        .wait()?;

    package_dir.push(package_name);

    Ok((status, package_dir))
}

fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
