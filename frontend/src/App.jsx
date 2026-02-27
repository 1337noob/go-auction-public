import React, { useState } from 'react'
import AuctionPage from './AuctionPage'

function App() {
  const [auctionId, setAuctionId] = useState('')
  const [userId, setUserId] = useState('')
  const [showAuction, setShowAuction] = useState(false)

  const handleStart = () => {
    if (auctionId && userId) {
      setShowAuction(true)
    } else {
      alert('Пожалуйста, введите ID аукциона и ID пользователя')
    }
  }

  if (showAuction) {
    return <AuctionPage auctionId={auctionId} userId={userId} />
  }

  return (
    <div className="app-container">
      <div className="login-form">
        <h1>Подключение к аукциону</h1>
        <div className="form-group">
          <label>ID Аукциона:</label>
          <input
            type="text"
            value={auctionId}
            onChange={(e) => setAuctionId(e.target.value)}
            placeholder="Введите ID аукциона"
          />
        </div>
        <div className="form-group">
          <label>ID Пользователя:</label>
          <input
            type="text"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            placeholder="Введите ID пользователя"
          />
        </div>
        <button onClick={handleStart}>Подключиться</button>
      </div>
    </div>
  )
}

export default App

