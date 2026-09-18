use std::process::Command;

use crate::{error::OurError, package_config::Config};

pub fn build(package_dir: &std::path::Path, config: &Config) -> Result<(), OurError> {
    let exit_status = Command::new(&config.build.command)
        .args(&config.build.args)
        .current_dir(package_dir)
        .spawn()?
        .wait()?;

    if !exit_status.success() {
        return Err(OurError::BuildFailed);
    }

    Ok(())
}
