use clap::{Parser, Subcommand};

mod build;
mod cargo;
mod error;
mod git;
mod package_config;
mod paths;
mod project;

#[derive(Parser)]
#[command()]
struct Args {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Install {
        url: String,

        #[arg(long)]
        branch: Option<String>,
        #[arg(short, long)]
        package_path: Option<String>,
    },
    Remove,
    Update,
}

fn main() {
    let args = Args::parse();

    match args.command {
        Commands::Install {
            url,
            branch,
            package_path,
        } => match install(url, branch, package_path) {
            Ok(()) => println!("successfully installed"),
            Err(error) => println!("{}", error),
        },
        Commands::Remove => println!("removing"),
        Commands::Update => println!("updating"),
    }
}

fn install(
    url: String,
    branch: Option<String>,
    package_path: Option<String>,
) -> Result<(), error::OurError> {
    println!("installing {}", url);

    let package_info = package_config::PackageInfo {
        name: git::get_name_from_url(&url).to_string(),
        url,
        branch,
        path: package_path,
    };

    let package_dir = git::clone(&package_info)?;
    println!("successfully cloned");

    let config = match package_config::load_config(package_dir.as_path())? {
        Some(config) => config,
        None => {
            let project_type = project::detect(&package_dir);
            package_config::generate_config(project_type, package_info, &package_dir)?
        }
    };

    build::build(&package_dir, &config)?;
    println!("successfully built");

    cargo::install_binary(&config)?;

    paths::clear_temp_dir()?;

    Ok(())
}
