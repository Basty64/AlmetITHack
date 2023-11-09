import './App.css';
import {useState} from "react";

function App() {
    const [message, setMessage] = useState('');

    const handleTextChange = (event) => {
        setMessage(event.target.value);
    };

    function handleSendButton() {
        const queryParam = JSON.stringify({message})
        console.log(queryParam)
        fetch('http://localhost:8080/article', {
        
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            // body: JSON.stringify({text}),
            body: queryParam
        })
            .then(response => response.json())
            .then(data => console.log(data))
            .catch(error => console.error('Error:', error));
    }

    return (
        <div>
            <textarea value={message} onChange={handleTextChange} />
            <button onClick={handleSendButton}>
            </button>
        </div>
    );
}

export default App;