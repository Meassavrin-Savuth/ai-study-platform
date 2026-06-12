import sys
import json
from youtube_transcript_api import YouTubeTranscriptApi
from youtube_transcript_api._errors import TranscriptsDisabled, NoTranscriptFound, VideoUnavailable

def get_video_id(url):
    if "v=" in url:
        return url.split("v=")[1].split("&")[0]
    if "youtu.be/" in url:
        return url.split("youtu.be/")[1].split("?")[0]
    return url

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No URL provided"}))
        sys.exit(1)

    url = sys.argv[1]
    video_id = get_video_id(url)

    try:
        ytt_api = YouTubeTranscriptApi()
        transcript = ytt_api.fetch(video_id)
        full_text = " ".join([line.text for line in transcript])
        print(json.dumps({"transcript": full_text}))

    except VideoUnavailable:
        print(json.dumps({"error": "This video is unavailable or private."}))

    except TranscriptsDisabled:
        print(json.dumps({"error": "This video does not have subtitles enabled. Please try a video with captions."}))

    except NoTranscriptFound:
        print(json.dumps({"error": "No subtitles found for this video. Please try a video with captions."}))

    except Exception as e:
        print(json.dumps({"error": "Could not extract transcript. Please try a different video."}))

main()