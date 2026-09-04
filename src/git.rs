use std::process::{Command, ExitStatus};

pub fn clone(url: &str) -> Result<ExitStatus, std::io::Error> {
    match Command::new("git").args(["clone", url]).spawn() {
        Ok(mut child) => match child.wait() {
            Ok(exit_status) => Ok(exit_status),
            Err(error) => Err(error),
        },
        Err(error) => Err(error),
    }
}
