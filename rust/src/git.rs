use crate::{error::OurError, paths::get_temp_dir};
use std::{path::PathBuf, process::Command};

pub fn clone(
    url: &str,
    branch: Option<String>,
    package_path: Option<String>,
) -> Result<PathBuf, OurError> {
    let package_name = get_name_from_url(url);
    let temp_dir = get_temp_dir()?;

    let mut args = vec!["clone".to_string(), "--depth".to_string(), "1".to_string()];

    if let Some(branch) = branch {
        args.extend(["--branch".to_string(), branch]);
    }

    args.extend([url.to_string(), package_name.to_string()]);

    let status = Command::new("git")
        .args(args)
        .current_dir(&temp_dir)
        .spawn()?
        .wait()?;

    if !status.success() {
        return Err(OurError::CloneFailed);
    }

    let mut package_dir = temp_dir;
    package_dir.push(package_name);

    if let Some(package_path) = package_path {
        package_dir.push(package_path);
    }

    Ok(package_dir)
}

fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
