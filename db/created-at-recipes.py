import json
import yt_dlp
from os import listdir
from os.path import join
import datetime

file_names = [f for f in listdir("./recipes-json") if f.endswith(".json")]

no_created_at = []

for file_name in file_names:
    file_path = join("./recipes-json", file_name)

    with open(file_path, "r", encoding="utf-8") as f:
        recipe = json.load(f)

    video_url = recipe.get("video_url")

    if video_url:
        print(recipe["title"])
        print(video_url)

        ydl_opts = {
            "cookiesfile": "cookies.txt",
            "verbose": False
        }

        try:
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                info = ydl.extract_info(video_url, download=False)

            timestamp = info.get("timestamp")
            if timestamp:
                recipe["created_at"] = datetime.datetime.fromtimestamp(
                    timestamp).strftime("%Y-%m-%d %H:%M:%S")
            else:
                recipe["created_at"] = ""

        except Exception as e:
            print(f"Failed to extract info for {video_url}: {e}")
            recipe["created_at"] = ""
            no_created_at.append(recipe["title"])
    else:
        recipe["created_at"] = ""
        no_created_at.append(recipe["title"])

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(recipe, f, ensure_ascii=False, indent=4)

print(no_created_at)
