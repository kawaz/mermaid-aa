//! mermaid-aa: Convert Mermaid diagrams to ASCII Art.

mod parser;
mod renderer;
mod types;
mod width;

use clap::Parser;
use std::io::{self, Read};
use std::path::PathBuf;
use std::process::ExitCode;

use parser::parse;
use renderer::{render, Charset, RenderOptions};
use types::Direction;
use width::AmbiguousWidth;

#[derive(Parser, Debug)]
#[command(name = "mermaid-aa")]
#[command(version, about = "Convert Mermaid diagrams to ASCII Art")]
struct Args {
    /// Mermaid text or file path
    #[arg()]
    input: Option<String>,

    /// Read from file
    #[arg(short, long)]
    file: Option<PathBuf>,

    /// Character set
    #[arg(short, long, default_value = "unicode", env = "MERMAID_AA_CHARSET")]
    charset: String,

    /// Width for East Asian Ambiguous characters
    ///
    /// Values:
    ///   1, half     Half-width (1 cell) - for Western terminals
    ///   2, full     Full-width (2 cells) - for CJK terminals
    ///   console     Ambiguous=1, Box Drawing=1 (recommended)
    ///   legacy      All ambiguous including Box Drawing=2
    #[arg(short, long, default_value = "1", env = "MERMAID_AA_AMBIGUOUS_WIDTH")]
    ambiguous_width: String,

    /// Override direction (TB, TD, BT, LR, RL)
    #[arg(short, long)]
    direction: Option<String>,
}

fn main() -> ExitCode {
    let args = Args::parse();

    // Parse charset
    let charset = match args.charset.parse::<Charset>() {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error: {}", e);
            return ExitCode::from(3);
        }
    };

    // Parse ambiguous width
    let ambiguous_width = match args.ambiguous_width.parse::<AmbiguousWidth>() {
        Ok(w) => w,
        Err(e) => {
            eprintln!("Error: {}", e);
            return ExitCode::from(3);
        }
    };

    // Parse direction override
    let direction_override = args
        .direction
        .as_ref()
        .and_then(|d| d.parse::<Direction>().ok());

    // Get input (priority: --file > argument > stdin)
    let input = match get_input(&args) {
        Ok(s) => s,
        Err(e) => {
            eprintln!("Error: {}", e);
            return ExitCode::from(2);
        }
    };

    // Parse mermaid diagram
    let flowchart = match parse(&input) {
        Ok(fc) => fc,
        Err(e) => {
            eprintln!("Error: {}", e);
            return ExitCode::from(1);
        }
    };

    // Render
    let options = RenderOptions {
        charset,
        ambiguous_width,
        direction_override,
    };

    let output = render(&flowchart, &options);
    println!("{}", output);

    ExitCode::SUCCESS
}

fn get_input(args: &Args) -> io::Result<String> {
    // Priority 1: --file option
    if let Some(path) = &args.file {
        return std::fs::read_to_string(path);
    }

    // Priority 2: positional argument
    if let Some(input) = &args.input {
        // Check if it's a file path
        let path = PathBuf::from(input);
        if path.exists() && path.is_file() {
            return std::fs::read_to_string(path);
        }
        // Otherwise treat as literal mermaid text
        return Ok(input.clone());
    }

    // Priority 3: stdin
    let mut buffer = String::new();
    io::stdin().read_to_string(&mut buffer)?;
    Ok(buffer)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_full_pipeline() {
        let input = r#"
flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[End]
    B -->|No| D[Retry]
"#;
        let fc = parse(input).unwrap();
        let output = render(&fc, &RenderOptions::default());

        assert!(output.contains("Start"));
        assert!(output.contains("Decision"));
        assert!(output.contains("End"));
        assert!(output.contains("Retry"));
    }

    #[test]
    fn test_japanese_labels() {
        let input = "flowchart TD\n    A[開始] --> B[終了]";
        let fc = parse(input).unwrap();

        // Test with half-width ambiguous
        let options = RenderOptions {
            ambiguous_width: AmbiguousWidth::Half,
            ..Default::default()
        };
        let output = render(&fc, &options);
        assert!(output.contains("開始"));
        assert!(output.contains("終了"));
    }

    #[test]
    fn test_all_charsets() {
        let input = "flowchart LR\n    A --> B";
        let fc = parse(input).unwrap();

        for charset in [
            Charset::Ascii,
            Charset::Unicode,
            Charset::UnicodeRound,
            Charset::UnicodeBold,
            Charset::UnicodeDouble,
        ] {
            let options = RenderOptions {
                charset,
                ..Default::default()
            };
            let output = render(&fc, &options);
            assert!(!output.is_empty());
        }
    }

    #[test]
    fn test_all_directions() {
        let input = "flowchart TD\n    A --> B";
        let fc = parse(input).unwrap();

        for dir in [Direction::TB, Direction::BT, Direction::LR, Direction::RL] {
            let options = RenderOptions {
                direction_override: Some(dir),
                ..Default::default()
            };
            let output = render(&fc, &options);
            assert!(!output.is_empty());
        }
    }
}
