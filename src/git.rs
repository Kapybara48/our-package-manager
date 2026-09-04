use std::process::Command;

pub fn clone(url: &str) {
    match Command::new("git").args(["clone", url]).spawn() {
        Ok(mut child) => match child.wait() {
            Ok(exit_status) => println!("Exit status: {}", exit_status),
            Err(error) => println!("{}", error),
        },
        Err(error) => println!("{}", error),
    }
}
