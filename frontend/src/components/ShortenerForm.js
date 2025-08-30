import React, { useState } from 'react';
import '../App.css';
import { FaHeart, FaGithub } from 'react-icons/fa';

const Shortener = () => {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!url.trim()) {
      setError('Please enter a URL');
      return;
    }

    setLoading(true);
    setError('');
    setShortUrl('');

    try {
      const response = await fetch('http://localhost:8080/shorten', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: url.trim() }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to shorten URL');
      }

      setShortUrl(data.shortUrl);
      setUrl('');
    } catch (err) {
      setError(err.message || 'Network error. Check if backend is running.');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(shortUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      const textArea = document.createElement('textarea');
      textArea.value = shortUrl;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="shortener-page">
      {/* === HEADER === */}
      <header className="shortener-header">
        <h1>
          <span className="gradient-text">Shorten Your Loooong Links :)</span>
        </h1>
        <p className="subtitle">
          Linky is an efficient and easy-to-use URL shortening service that streamlines your online experience.
        </p>
      </header>

      {/* === FORM === */}
      <main className="shortener-main">
        <form onSubmit={handleSubmit} className="shortener-form">
          <div className="input-wrapper">
            <input
              type="text"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Enter the link here"
              disabled={loading}
              className="url-input"
            />
            <button
              type="submit"
              disabled={loading || !url.trim()}
              className="shorten-btn"
            >
              {loading ? 'Working...' : 'Shorten Now!'}
            </button>
          </div>
        </form>

        {error && <div className="error-msg">⚠️ {error}</div>}

        {shortUrl && (
          <div className="result-bar">
            <input type="text" value={shortUrl} readOnly className="result-input" />
            <button onClick={copyToClipboard} className="copy-btn">
              {copied ? '✅ Copied!' : '📋 Copy'}
            </button>
          </div>
        )}

        <div className="footer-note">
          <p>
            You can create <span className="highlight">05</span> more links.{' '}
          </p>
        </div>
      </main>

      {/* === NEW FOOTER SECTION === */}
      <footer className="shortener-footer">
        <p>
          Made with <FaHeart className="icon-heart" /> by{' '}
          <a href="https://github.com/sachinky09" target="_blank" rel="noreferrer">
            <FaGithub className="icon-github" /> sachinky09
          </a>
        </p>
        <p>
          🛠️ <a href="https://github.com/sachinky09/url-shortener" target="_blank" rel="noreferrer">
            Clone and customise – click here
          </a>
        </p>
      </footer>
    </div>
  );
};

export default Shortener;
