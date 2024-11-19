AI Name Validator
=====================

This Go program was made for fun and just to explore using the Ollama API. It processes a CSV file containing user data and determines if each user's full name is likely a person's name using a local language model API.

Features
--------

*   Reads a CSV file (`users.csv`) containing user data.
*   Extracts `firstName` and `lastName` columns to form full names.
*   Queries a language model API to check if the full name is likely a person's name.
*   Displays processing progress with a spinner and percentage complete.
*   Logs names that are not likely to be a person's name.

Prerequisites
-------------

*   **Go** installed (version 1.16 or higher).
*   **Ollama Server** running and accessible at `http://localhost:11434/api/generate`.

Installation
------------

1.  **Clone the repository**:
    
    ```bash
    git clone nathan-power/ai-name-validator
    cd nathan-power/ai-name-validator
    ```
    
2.  **Prepare the CSV file**:
    
    Ensure you have a CSV file named `users.csv` in the same directory, with at least the following headers:
    
    *   `firstName`
    *   `lastName`
    
    Example `users.csv`:
    
    ```csv
    firstName,lastName,email
    John,Doe,john.doe@example.com
    Jane,Smith,jane.smith@example.com
    ```
    
3.  **Set up the Ollama Server**:
    
    The program requires a language model API that accepts POST requests and responds according to the specified format.
    
    **API Endpoint**: `http://localhost:11434/api/generate`
    
    **Request Format**:
    
    *   Method: `POST`
        
    *   Headers: `Content-Type: application/json`
        
    *   Body:
        
        ```json
        {
          "model": "phi3",
          "prompt": "<Prompt string>",
          "temperature": 0
        }
        ```
        
    
    **Response Format**:
    
    *   Status Code: `200 OK`
        
    *   Body:
        
        ```json
        {
          "response": "<yes|no>"
        }
        ```
        
    
    **Example Setup**:
    
    You can implement a simple API server using any language or framework that you prefer. The server should parse the incoming JSON, process the `prompt`, and return a JSON response with a `response` field containing `yes` or `no`.
    
    _Note_: The implementation of the language model API server is beyond the scope of this README. Ensure that it meets the above specifications.
    

Usage
-----

1.  **Build the program**:
    
    ```bash
    go build -o ai-name-validator
    ```
    
2.  **Run the program**:
    
    ```bash
    ./ai-name-validator
    ```
    
    The program will process the `users.csv` file and output names that are not likely a person's name.
    

Example Output
--------------

```vbnet
Processing... 25% complete /
John Doe is not likely a person's name.
Processing... 50% complete -
Processing... 75% complete \
Jane Smith is not likely a person's name.
Processing... 100% complete |
Processing complete.
```

Customization
-------------

*   **CSV File**: To process a different CSV file, modify the `csvFileName` variable in the `processRecords` function call within the `main` function.
    
*   **API Endpoint**: If your language model API is running on a different URL or port, update the `url` variable in the `postData` function.
    

Troubleshooting
---------------

*   **CSV File Errors**:
    
    *   Ensure that the `users.csv` file exists in the working directory.
    *   Verify that the CSV file has the required headers: `firstName`, `lastName`.
*   **API Connection Errors**:
    
    *   Ensure that the language model API server is running and accessible.
    *   Verify that the server is listening on `http://localhost:11434/api/generate`.
*   **JSON Decoding Errors**:
    
    *   Check that the API server responds with the expected JSON format.
    *   Ensure that the `response` field in the API response contains only `yes` or `no`.

Code Structure
--------------

*   **main.go**: The main Go file containing all the logic for processing the CSV and interacting with the API.
    
*   **Functions**:
    
    *   `main()`: Entry point of the program. Calls `processRecords`.
    *   `processRecords(csvFileName string)`: Opens the CSV file, reads headers, and initiates record processing.
    *   `processEachRecord(...)`: Processes each CSV record, constructs full names, and queries the model.
    *   `queryModel(fullName string)`: Prepares the request to the language model API and retrieves the response.
    *   `postData(data []byte)`: Sends the HTTP POST request to the API and decodes the response.
    *   `lineCounter(fileName string)`: Counts the total number of lines in the CSV file.
    *   `findIndex(headers []string, column string)`: Finds the index of a specific column in the CSV headers.
    *   `displayProgress(...)`: Displays the processing progress with a spinner and percentage.
