use clap::{Parser, Subcommand};

mod build;
mod cargo;
mod error;
mod git;
mod package_config;
mod paths;

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
        } => match install(&url, branch, package_path) {
            Ok(()) => println!("successfully installed"),
            Err(error) => println!("{}", error),
        },
        Commands::Remove => println!("removing"),
        Commands::Update => println!("updating"),
    }
}

fn install(
    url: &str,
    branch: Option<String>,
    package_path: Option<String>,
) -> Result<(), error::OurError> {
    println!("installing {}", url);

    let package_dir = git::clone(url, branch, package_path)?;
    println!("successfully cloned");

    build::build(&package_dir)?;
    println!("successfully built");

    cargo::install_binary(&package_dir)?;

    Ok(())
}
