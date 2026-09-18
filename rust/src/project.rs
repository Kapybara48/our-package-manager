use std::path::Path;

pub enum ProjectType {
    Rust,
    Unknown,
}

pub fn detect(package_dir: &Path) -> ProjectType {
    let mut cargo_file = package_dir.to_path_buf();
    cargo_file.push("Cargo.toml");

    if cargo_file.exists() {
        ProjectType::Rust
    } else {
        ProjectType::Unknown
    }
}
