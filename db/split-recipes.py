import re
from pathlib import Path

recipe_names = []

input_file = "./QTcordCookbook.md"          # your big markdown file
output_dir = Path("recipes")  # folder for outputs
output_dir.mkdir(exist_ok=True)

# Read the whole markdown
text = Path(input_file).read_text(encoding="utf-8")

# Split at each "## **Title**"
# (?=...) keeps the delimiter as part of each section
sections = re.split(r'(?=^##\s+\*\*.+?\*\*)', text, flags=re.MULTILINE)

# Remove empty pieces
sections = [s.strip() for s in sections if s.strip()]


def safe_filename(title: str) -> str:
    # Remove markdown stars and illegal filename chars
    title = re.sub(r'[*\\/:?"<>|]', '', title)
    title = title.strip().replace(" ", "_")

    if title in recipe_names:
        title = title + "-duplicate"

    recipe_names.append(title)
    return title


for section in sections:
    # Extract title from first line
    first_line = section.splitlines()[0]
    title_match = re.search(r'\*\*(.+?)\*\*', first_line)

    if title_match:
        name = safe_filename(title_match.group(1))
    else:
        name = "untitled"

    out_path = output_dir / f"{name}.md"
    out_path.write_text(section + "\n", encoding="utf-8")

print(f"Saved {len(sections)} files to {output_dir}")
