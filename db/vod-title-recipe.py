import json
import yt_dlp
from os import listdir
from os.path import join

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

    video_url = recipe.get("video_url")

    if video_url:
        print(recipe["title"])
        print(video_url)

        ydl_opts = {
            "verbose": False
        }

        try:
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                info = ydl.extract_info(video_url, download=False)

                recipe["vod_title"] = info["title"]

        except Exception as e:
            print(f"Failed to extract info for {video_url}: {e}")
    else:
        recipe["vod_title"] = ""

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)
