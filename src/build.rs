use std::process::{Command, ExitStatus};

pub fn build() -> Result<ExitStatus, std::io::Error> {
    match Command::new("cargo").arg("build").spawn() {
        Ok(mut child) => match child.wait() {
            Ok(exit_status) => Ok(exit_status),
            Err(error) => Err(error),
        },
        Err(error) => Err(error),
    }
}
