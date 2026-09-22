use crate::{error::OurError, package_config::PackageInfo, paths::get_temp_dir};
use std::{path::PathBuf, process::Command};

pub fn clone(package_info: &PackageInfo) -> Result<PathBuf, OurError> {
    let temp_dir = get_temp_dir()?;

    let mut args = vec!["clone".to_string(), "--depth".to_string(), "1".to_string()];

    if let Some(branch) = &package_info.branch {
        args.extend(["--branch".to_string(), branch.clone()]);
    }

    args.extend([package_info.url.clone(), package_info.name.clone()]);

    let status = Command::new("git")
        .args(args)
        .current_dir(&temp_dir)
        .spawn()?
        .wait()?;

    if !status.success() {
        return Err(OurError::CloneFailed);
    }

    let mut package_dir = temp_dir;
    package_dir.push(&package_info.name);

    if let Some(package_path) = &package_info.path {
        package_dir.push(package_path);
    }

    Ok(package_dir)
}

pub fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
