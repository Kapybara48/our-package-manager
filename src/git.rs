use std::process::{Command, ExitStatus};

pub fn clone(url: &str) -> Result<ExitStatus, std::io::Error> {
    match Command::new("git")
        .args(["clone", url, get_name_from_url(url)])
        .spawn()
    {
        Ok(mut child) => match child.wait() {
            Ok(exit_status) => Ok(exit_status),
            Err(error) => Err(error),
        },
        Err(error) => Err(error),
    }
}

fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
