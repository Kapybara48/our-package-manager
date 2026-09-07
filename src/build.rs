use std::process::{Command, ExitStatus};

use crate::error::OurError;

pub fn build(package_dir: std::path::PathBuf) -> Result<ExitStatus, OurError> {
    Ok(Command::new("cargo")
        .arg("build")
        .current_dir(package_dir)
        .spawn()?
        .wait()?)
}
