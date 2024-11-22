            Welcome to the ASCII Art Generator website! This website allows you to turn text into ASCII art using
            different styles. You can download the content in text or HTML formats as well.
  

            **Description**
            This project, as part of the ZONE01 coursework, comprises several exercises aimed at mastering various
            aspects of web development, particularly using Go and CSS. Our focus has been on creating a user-friendly
            interface, implementing robust functionalities, and ensuring the overall consistency and responsiveness of
            our web application.

            **Technologies**
            - GOLANG
            - HTML
            - CSS
            - JavaScript
            
            **Installation**
            To install and run the project locally:
            1. Clone the repository:
            git clone https://github.com/OsBTD/ascii-art-web-complete.git
        
            2. Navigate to the project directory:
            cd ascii-art-web-complete
           
            3. Build and run the Go server:
            go run main.go
       
            4. Open your web browser and visit `http://localhost:8080`.
  
            **To run the project using Docker:**
            - Ensure Docker is installed and running on your machine.
            - Pull the Docker image:
            docker pull ascii-art-web
            - Run the Docker container:
            docker run -p 8080:8080 ascii-art-web
    
            **Docker Commands Explained:**
            - docker build -t ascii-art-web . : Builds the Docker image from the Dockerfile in the
            current directory and tags it as ascii-art-web.
            - docker run -p 8080:8080 ascii-art-web : Runs a container from the ascii-art-web image,
            mapping port 8080 on your machine to port 8080 in the container.
    
            If you make changes to your code, rebuild the image using:
            docker build -t ascii-art-web .
         
    
                **Usage**
                - **Text Input**: Enter the text you want to convert into ASCII art.
                - **Banner Selection**: Choose from various banners.
                - **Generate**: Click the button to generate the ASCII art.
                - **Export**: Download the ASCII art in text or HTML formats or copy it to the clipboard.
    
                **Implementation Details**
                - **Algorithm**: The ASCII art generation algorithm uses predefined banners to convert text characters into
                corresponding ASCII art representations.
                - **HTTP Server**: The server is written in Go and handles requests for generating and exporting ASCII art.
                - **Docker**: The project includes a Dockerfile for containerization, following best practices for creating
                images and managing containers.
    
                **Join Us on Our Journey**
                We welcome contributions! Please fork the repository and submit pull requests.
    
                **Contact**
                If you have any questions or feedback, feel free to contact us at our email: borntolearn314@gmail.com.

