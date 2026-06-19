import sys
import json
import os
import traceback
from youtube_transcript_api import YouTubeTranscriptApi
from youtube_transcript_api._errors import TranscriptsDisabled, NoTranscriptFound, VideoUnavailable

def get_video_id(url):
    """Extract video ID from various YouTube URL formats."""
    url = url.strip()
    if "youtu.be/" in url:
        return url.split("youtu.be/")[1].split("?")[0].split("&")[0]
    if "v=" in url:
        return url.split("v=")[1].split("&")[0].split("?")[0]
    if "/shorts/" in url:
        return url.split("/shorts/")[1].split("?")[0].split("&")[0]
    if "/embed/" in url:
        return url.split("/embed/")[1].split("?")[0].split("&")[0]
    return url

def make_api(cookies_path=None):
    """Create YouTubeTranscriptApi instance, with cookies if available."""
    if cookies_path and os.path.exists(cookies_path):
        return YouTubeTranscriptApi(cookie_path=cookies_path)
    return YouTubeTranscriptApi()

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No URL provided"}))
        sys.exit(1)

    url = sys.argv[1]
    video_id = get_video_id(url)

    # Look for cookies.txt in the scripts directory or backend root
    script_dir = os.path.dirname(os.path.abspath(__file__))
    cookies_candidates = [
        os.path.join(script_dir, "cookies.txt"),
        os.path.join(script_dir, "..", "cookies.txt"),
    ]
    cookies_path = next((p for p in cookies_candidates if os.path.exists(p)), None)

    try:
        ytt_api = make_api(cookies_path)

        # First try: fetch directly (no language preference)
        try:
            transcript = ytt_api.fetch(video_id)
            full_text = " ".join([line.text for line in transcript])
            print(json.dumps({"transcript": full_text}))
            return
        except Exception:
            pass

        # Second try: list all and pick best available
        try:
            transcript_list = ytt_api.list(video_id)
            transcripts = list(transcript_list)

            if not transcripts:
                print(json.dumps({"error": "No subtitles found for this video. Please try a video with captions."}))
                return

            # prefer manually created over auto-generated
            chosen = next((t for t in transcripts if not t.is_generated), transcripts[0])
            fetched = chosen.fetch()
            full_text = " ".join([line.text for line in fetched])
            print(json.dumps({"transcript": full_text}))
            return

        except Exception as inner_e:
            print(json.dumps({"error": f"Could not fetch transcript: {str(inner_e)}"}))
            return

    except VideoUnavailable:
        print(json.dumps({"error": "This video is unavailable or private."}))

    except TranscriptsDisabled:
        print(json.dumps({"error": "This video does not have subtitles enabled. Please try a video with captions."}))

    except NoTranscriptFound:
        print(json.dumps({"error": "No subtitles found for this video. Please try a video with captions."}))

    except Exception as e:
        print(json.dumps({"error": f"Unexpected error: {str(e)}", "trace": traceback.format_exc()}))

main()
