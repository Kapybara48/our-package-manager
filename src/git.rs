use crate::{error::OurError, paths::get_our_dir};
use std::process::{Command, ExitStatus};

pub fn clone(url: &str) -> Result<ExitStatus, OurError> {
    Ok(Command::new("git")
        .args(["clone", url, get_name_from_url(url)])
        .current_dir(get_our_dir()?)
        .spawn()?
        .wait()?)
}

fn get_name_from_url(url: &str) -> &str {
    url.trim_end_matches(".git")
        .rsplit("/")
        .next()
        .unwrap_or("package")
}
