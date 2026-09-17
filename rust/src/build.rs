use std::process::Command;

use crate::error::OurError;

pub fn build(package_dir: &std::path::Path) -> Result<(), OurError> {
    let exit_status = Command::new("cargo")
        .arg("build")
        .current_dir(package_dir)
        .spawn()?
        .wait()?;

    if !exit_status.success() {
        return Err(OurError::BuildFailed);
    }

    Ok(())
}
