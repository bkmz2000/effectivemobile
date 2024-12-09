import requests
import os
import time

# Load environment variables
DB_HOST = os.getenv('DB_HOST', 'localhost')
DB_PORT = os.getenv('DB_PORT', '5432')
DB_USER = os.getenv('DB_USER', 'your_db_user')
DB_PASSWORD = os.getenv('DB_PASSWORD', 'your_db_password')
DB_NAME = os.getenv('DB_NAME', 'your_db_name')

# API endpoint
API_URL = "http://localhost:8080/song"

def populate_db():
    # Create sample song entries using the API
    songs = [
        {
            "name": "Test Song",
            "group": "Test Group",
            "release_date": "2023-01-01",
            "text": "This is a test song.",
            "link": "http://example.com"
        },
        {
            "name": "Another Song",
            "group": "Test Group",
            "release_date": "2023-01-02",
            "text": "This is another test song.",
            "link": "http://example.com/another"
        },
        {
            "name": "Third Song",
            "group": "Test Group",
            "release_date": "2023-01-03",
            "text": "This is the third test song.",
            "link": "http://example.com/third"
        },
        {
            "name": "Fourth Song",
            "group": "Different Group",
            "release_date": "2023-01-04",
            "text": "This is the fourth test song.",
            "link": "http://example.com/fourth"
        },
        {
            "name": "Fifth Song",
            "group": "Test Group",
            "release_date": "2023-01-05",
            "text": "This is the fifth test song.",
            "link": "http://example.com/fifth"
        },
        {
            "name": "Sixth Song",
            "group": "Different Group",
            "release_date": "2023-01-06",
            "text": "This is the sixth test song.",
            "link": "http://example.com/sixth"
        }
    ]

    for song in songs:
        response = requests.post(API_URL, json=song)
        print(f"Added song: {song['name']}")

def test_service():
    # Test GET request
    response = requests.get(API_URL, params={"group": "Test Group", "name": "Test Song"})
    print("GET Response:", response.json())

    # Test POST request
    new_song = {
        "name": "New Test Song",
        "group": "New Test Group",
        "release_date": "2023-01-01",
        "text": "This is a new test song.",
        "link": "http://example.com/new"
    }
    response = requests.post(API_URL, json=new_song)
    print("POST Response:", response.json())


def test_filtering():
    # Test filtering by group
    response = requests.get(API_URL, params={"group": "Test Group"})
    print("Filter by Group Response:", response.json())

    # Test filtering by name
    response = requests.get(API_URL, params={"name": "Test Song"})
    print("Filter by Name Response:", response.json())

    # Test filtering by both group and name
    response = requests.get(API_URL, params={"group": "Test Group", "name": "Another Song"})
    print("Filter by Group and Name Response:", response.json())

    # Test filtering by date range
    response = requests.get(API_URL, params={"group": "Test Group", "date_after": "2023-01-01", "date_before": "2023-01-03"})
    print("Filter by Date Range Response:", response.json())

def clean_db():
    # Clean up the test data using DELETE method
    songs_to_delete = ['Test Song', 'New Test Song', 'Another Song', 'Third Song', 'Fifth Song', 'Fourth Song', 'Sixth Song']
    
    for song_name in songs_to_delete:
        print(f"deleting {song_name}")
        response = requests.get(API_URL, params={"name": song_name})
        
        if response.status_code == 200:
            song_data = response.json()
            song_id = song_data[0].get('id')  # Assuming the response contains the ID of the song
            
            if song_id:
                # Now delete the song by ID
                delete_response = requests.delete(API_URL, params={"id": song_id})
                print(f"Clean DB Response for {song_name} (ID: {song_id}):", delete_response.status_code)
            else:
                print(f"No ID found for song: {song_name}")
        else:
            print(f"Failed to retrieve song: {song_name}, Status Code: {response.status_code}")

if __name__ == "__main__":
    populate_db()  # Use the API to populate the database
    test_service()
    test_filtering()  # Call the filtering test
    clean_db()