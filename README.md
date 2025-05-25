# FF7Book: Final Fantasy VII AI Novel Generator

A command-line tool written in Go that converts the _Final Fantasy VII_ game script and community story summary into a polished novel.
This fan project streams AI-generated prose—formatted in Markdown—ready for eBook or print conversion.

## Example output

A sample output of the generated novel can be found in the `data/output-20250525.md` file.
The novel is incomplete because the process reached the token limit of the AI model.

[FF7Book Example Output](data/output-20250525.md)

## Installation

```bash
go install github.com/igolaizola/ff7book/cmd/ff7book@latest
```

## Usage

You can supply options via flags, environment variables (prefixed `FF7BOOK_`), or a YAML config file.

```yaml
# example.yaml
key: "YOUR_GOOGLE_API_KEY"
input: "data/ff7_script.txt"
output: "ff7_novel.md"
model: "gemini-2.5-pro-preview-05-06"
debug: false
```

```bash
ff7book generate \
  --key "$GOOGLE_API_KEY" \
  --input data/ff7_script.txt \
  --output ff7_novel.md \
  --model gemini-2.5-pro-preview-05-06
```

Or using config file:

```bash
ff7book generate --config example.yaml
```

**Flags:**

| Flag       | Description                               |
| ---------- | ----------------------------------------- |
| `--key`    | AI provider API key                       |
| `--input`  | Path to the game script + synopsis        |
| `--output` | Destination file for the novel (Markdown) |
| `--model`  | AI model name (e.g., Gemini, GPT-4)       |
| `--debug`  | Enable debug logging                      |

Example:

```bash
./ff7book generate --key "$GOOGLE_API_KEY" --input data/input.md --output ff7-novel.md
```

## Resources

- https://www.yinza.com/Fandom/Script.html
- https://finalfantasy.fandom.com/wiki/Final_Fantasy_VII

## Legal Notice & Disclaimer

**Unlicensed Fan Project**: This tool and its outputs are not endorsed by, affiliated with, or sponsored by Square Enix Co., Ltd.

**Intellectual Property**: _Final Fantasy VII_ characters, narratives, and trademarks belong to Square Enix. All source text is scraped from publicly available Fandom pages.

**Use Restrictions**: For **personal, educational, and non-commercial** use only.

**Fair Use**: Output is transformative—novelization for commentary and fan art. Rights holders may contact the repo owner for concerns.
