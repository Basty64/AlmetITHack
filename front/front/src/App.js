import './App.css';
import { useState } from "react";

function App() {
  const [message, setMessage] = useState('');
  const [responseText, setResponseText] = useState('');
  const [responseTags, setResponseTags] = useState('');
  const [articles, setArticles] = useState(null);

  const handleTextChange = (event) => {
    setMessage(event.target.value);
  };

  function handleSendButton() {
    const queryParam = JSON.stringify({ message });
    console.log(queryParam);
    fetch('http://localhost:8080/article', {
      // mode: 'no-cors',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: queryParam
    }).then(response => response.json())
      .then(data => {
        console.log(data);
        setResponseText(data.text);
        const concatenatedString = data.tags.reduce((accumulator, currentValue) => accumulator + "\n" + currentValue, '');
        setResponseTags(concatenatedString);
        setArticles(data.links[0].message);
      })
      .catch(error => console.error('Pizda:', error));
  }

  function renderArticles(articles) {
    if (!articles) {
      return null;
    }
    console.log(articles);
    return (
      <div>
        <h2>Articles</h2>
        {articles.items.map((item, index) => (
          <div key={index}>
            <h3>{item.title.join(', ')}</h3>
            <a href="{item.URL}" className="article-link">Ссылка на статью</a>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div>
      <div className="input-container">
        <textarea value={message} onChange={handleTextChange} className="message-input" />
        <button id="greateButton" onClick={handleSendButton} className="send-button">поиск</button>
      </div>
      <div className="response-container">
        <p className="response-text">{responseText}</p>
        <p className="response-tags">{responseTags}</p>
        <p className="links-heading">Links</p>
        {renderArticles(articles)}
      </div>
    </div>
  );
}

export default App;
