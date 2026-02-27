import React, { useState, useEffect, useRef } from 'react'

const API_BASE = '/api'
const WS_BASE = 'ws://localhost:8081'

function AuctionPage({ auctionId, userId }) {
  const [auction, setAuction] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [bidAmount, setBidAmount] = useState('')
  const [placingBid, setPlacingBid] = useState(false)
  const [bidError, setBidError] = useState(null)
  const [cancellationReason, setCancellationReason] = useState(null)
  const [timeoutEnd, setTimeoutEnd] = useState(null)
  const [timeRemaining, setTimeRemaining] = useState(null)
  const [auctionEndTime, setAuctionEndTime] = useState(null)
  const [auctionEndRemaining, setAuctionEndRemaining] = useState(null)
  const wsRef = useRef(null)
  const timeoutIntervalRef = useRef(null)
  const auctionEndIntervalRef = useRef(null)
  const shouldReconnectRef = useRef(true)
  const reconnectTimeoutRef = useRef(null)

  // Загрузка информации об аукционе
  useEffect(() => {
    loadAuction()
  }, [auctionId])

  // Подключение к WebSocket
  useEffect(() => {
    shouldReconnectRef.current = true
    connectWebSocket()
    return () => {
      // Отключаем переподключение при размонтировании компонента
      shouldReconnectRef.current = false
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
      }
      if (timeoutIntervalRef.current) {
        clearInterval(timeoutIntervalRef.current)
      }
      if (auctionEndIntervalRef.current) {
        clearInterval(auctionEndIntervalRef.current)
      }
    }
  }, [auctionId, userId])

  // Таймер обратного отсчета таймаута ставки
  useEffect(() => {
    if (timeoutEnd) {
      timeoutIntervalRef.current = setInterval(() => {
        const now = new Date().getTime()
        const remaining = Math.max(0, timeoutEnd - now)
        setTimeRemaining(remaining)

        if (remaining === 0) {
          clearInterval(timeoutIntervalRef.current)
          setTimeoutEnd(null)
        }
      }, 100)
    } else {
      if (timeoutIntervalRef.current) {
        clearInterval(timeoutIntervalRef.current)
      }
      setTimeRemaining(null)
    }

    return () => {
      if (timeoutIntervalRef.current) {
        clearInterval(timeoutIntervalRef.current)
      }
    }
  }, [timeoutEnd])

  // Таймер обратного отсчета до завершения аукциона
  useEffect(() => {
    if (auctionEndTime) {
      auctionEndIntervalRef.current = setInterval(() => {
        const now = new Date().getTime()
        const remaining = Math.max(0, auctionEndTime - now)
        setAuctionEndRemaining(remaining)

        if (remaining === 0) {
          clearInterval(auctionEndIntervalRef.current)
          setAuctionEndTime(null)
        }
      }, 100)
    } else {
      if (auctionEndIntervalRef.current) {
        clearInterval(auctionEndIntervalRef.current)
      }
      setAuctionEndRemaining(null)
    }

    return () => {
      if (auctionEndIntervalRef.current) {
        clearInterval(auctionEndIntervalRef.current)
      }
    }
  }, [auctionEndTime])

  const loadAuction = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${API_BASE}/auctions/${auctionId}`)
      if (!response.ok) {
        throw new Error('Не удалось загрузить информацию об аукционе')
      }
      const data = await response.json()
      setAuction(data)
      setError(null)

      // Отсчет работает только после старта аукциона
      if (data.status === 'started') {
        const now = new Date().getTime()
        
        if (data.current_bid) {
          // Если есть текущая ставка, устанавливаем таймаут
          const bidTime = parseTimeMs(data.current_bid.created_at)
          const timeoutMs = data.timeout ? parseTimeout(data.timeout) : 0
          const timeoutEndTime = bidTime + timeoutMs
          
          // Устанавливаем таймаут только если он еще не истек
          if (timeoutEndTime > now) {
            setTimeoutEnd(timeoutEndTime)
            // Обновляем время завершения: время ставки + timeout
            setAuctionEndTime(timeoutEndTime)
          } else {
            // Таймаут уже истек, не показываем его
            setTimeoutEnd(null)
            // Время завершения = текущее время + timeout (если ставок не будет)
            const newEndTime = now + timeoutMs
            const endTime = new Date(data.end_time).getTime()
            setAuctionEndTime(Math.min(endTime, newEndTime))
          }
        } else {
          // Если ставок нет, но аукцион начался, используем EndTime или время старта + timeout
          const endTime = parseTimeMs(data.end_time)
          const startTime = parseTimeMs(data.start_time)
          const timeoutMs = data.timeout ? parseTimeout(data.timeout) : 0
          
          // Время завершения = min(EndTime, StartTime + timeout)
          // Это время, когда аукцион завершится, если ставок не будет
          const calculatedEndTime = Math.min(endTime, startTime + timeoutMs)
          if (calculatedEndTime > now) {
            setAuctionEndTime(calculatedEndTime)
          }
          setTimeoutEnd(null)
        }
      } else {
        // Если аукцион еще не начался, не устанавливаем отсчет
        setTimeoutEnd(null)
        setAuctionEndTime(null)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const parseTimeMs = (value) => {
    if (!value) return null
    if (typeof value === 'number') return value

    const str = String(value).trim()

    // Определяем, есть ли таймзона
    const hasTZ = /[zZ]|[+-]\d{2}:?\d{2}$/.test(str)

    // Нормализуем ISO с лишними наносекундами.
    // Если таймзона отсутствует — считаем, что это UTC и добавляем 'Z'.
    const match = str.match(
      /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(\.\d+)?([+-]\d{2}:?\d{2}|Z)?$/
    )
    if (match) {
      const base = match[1]
      const frac = match[2] || ''
      const tz = match[3] || (hasTZ ? '' : 'Z')
      let msPart = ''
      if (frac) {
        // Оставляем только миллисекунды (3 знака)
        const trimmed = frac.slice(1, 4) // пропускаем точку
        msPart = `.${trimmed.padEnd(3, '0')}`
      }
      return new Date(`${base}${msPart}${tz}`).getTime()
    }

    // Fallback — отдаем в Date как есть (на случай нестандартных форматов)
    return new Date(hasTZ ? str : `${str}Z`).getTime()
  }

  const parseTimeout = (timeout) => {
    // В Go time.Duration сериализуется в JSON как наносекунды (число)
    if (typeof timeout === 'number') {
      return timeout / 1000000 // конвертируем наносекунды в миллисекунды
    }
    // Если строка, парсим формат "30s", "1m", "5m30s" и т.д.
    if (typeof timeout === 'string') {
      const match = timeout.match(/(\d+h)?(\d+m)?(\d+s)?/)
      if (!match) return 0
      
      let ms = 0
      if (match[1]) ms += parseInt(match[1]) * 3600000
      if (match[2]) ms += parseInt(match[2]) * 60000
      if (match[3]) ms += parseInt(match[3]) * 1000
      return ms
    }
    return 0
  }

  const connectWebSocket = () => {
    const ws = new WebSocket(`${WS_BASE}/ws?auction_id=${auctionId}&user_id=${userId}`)
    
    ws.onopen = () => {
      console.log('WebSocket подключен')
      // Очищаем таймер переподключения при успешном подключении
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
        reconnectTimeoutRef.current = null
      }
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        handleWebSocketEvent(message)
      } catch (err) {
        console.error('Ошибка парсинга сообщения:', err)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket ошибка:', error)
    }

    ws.onclose = (event) => {
      console.log('WebSocket отключен', event.code, event.reason)
      wsRef.current = null
      
      // Переподключение только если это не было намеренное закрытие
      // и компонент еще монтирован
      if (shouldReconnectRef.current && auctionId && userId) {
        console.log('Переподключение через 3 секунды...')
        reconnectTimeoutRef.current = setTimeout(() => {
          if (shouldReconnectRef.current && auctionId && userId) {
            connectWebSocket()
          }
        }, 3000)
      }
    }

    wsRef.current = ws
  }

  const handleWebSocketEvent = (message) => {
    const { event_type, event_data } = message

    switch (event_type) {
      case 'AuctionStarted':
        setAuction((prev) => {
          const updated = {
            ...prev,
            status: 'started'
          }
          
          // Если ставок нет, устанавливаем время завершения
          if (!prev.current_bid && prev.timeout) {
            const now = new Date().getTime()
            const timeoutMs = parseTimeout(prev.timeout)
            const endTime = parseTimeMs(prev.end_time)
            const startTime = parseTimeMs(prev.start_time) || now
            // Используем минимальное из EndTime и (время старта + timeout)
            const calculatedEndTime = Math.min(endTime, startTime + timeoutMs)
            if (calculatedEndTime > now) {
              setAuctionEndTime(calculatedEndTime)
            } else {
              // Если уже истекло, сбрасываем
              setAuctionEndTime(null)
            }
          }
          
          return updated
        })
        break

      case 'BidPlaced':
        if (event_data.auction_id === auctionId) {
          // Очищаем ошибку при успешной ставке
          if (event_data.user_id === userId) {
            setBidError(null)
          }
          
          // Устанавливаем новый таймаут ставки и время завершения
          const bidTime = parseTimeMs(event_data.timestamp)
          
          setAuction((prev) => {
            const timeoutMs = prev?.timeout ? parseTimeout(prev.timeout) : 0
            const newEndTime = bidTime + timeoutMs
            
            // Устанавливаем таймаут ставки
            setTimeoutEnd(newEndTime)
            
            // Обновляем время завершения: текущее время ставки + timeout
            setAuctionEndTime(newEndTime)
            
            // Обновляем текущую ставку
            return {
              ...prev,
              current_bid: {
                id: event_data.bid_id,
                user_id: event_data.user_id,
                amount: event_data.amount,
                created_at: event_data.timestamp
              }
            }
          })
        }
        break

      case 'BidRejected':
        if (event_data.auction_id === auctionId && event_data.user_id === userId) {
          // Показываем ошибку пользователю, если это его ставка была отклонена
          const errorMessage = event_data.error || 'Ставка была отклонена'
          setBidError(errorMessage)
          // Автоматически скрываем ошибку через 5 секунд
          setTimeout(() => setBidError(null), 5000)
        }
        break

      case 'AuctionCancelled':
        if (event_data.auction_id === auctionId) {
          setAuction((prev) => ({
            ...prev,
            status: 'cancelled'
          }))
          // Сохраняем причину отмены
          setCancellationReason(event_data.reason || 'Причина не указана')
          // Останавливаем все таймеры
          setTimeoutEnd(null)
          setAuctionEndTime(null)
        }
        break

      case 'AuctionCompleted':
        if (event_data.auction_id === auctionId) {
          setAuction((prev) => ({
            ...prev,
            status: 'completed',
            winner_id: event_data.winner_id,
            final_price: event_data.final_price
          }))
          setTimeoutEnd(null)
          setAuctionEndTime(null)
        }
        break

      default:
        console.log('Неизвестное событие:', event_type)
    }
  }

  const handlePlaceBid = async (e) => {
    e.preventDefault()
    
    if (!bidAmount || isNaN(bidAmount) || parseInt(bidAmount) <= 0) {
      alert('Введите корректную сумму ставки')
      return
    }

    setPlacingBid(true)
    try {
      const response = await fetch(`${API_BASE}/auctions/${auctionId}/bids`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          user_id: userId,
          amount: parseInt(bidAmount)
        })
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || 'Не удалось сделать ставку')
      }

      setBidAmount('')
    } catch (err) {
      alert(`Ошибка: ${err.message}`)
    } finally {
      setPlacingBid(false)
    }
  }

  const formatTime = (ms) => {
    if (!ms) return '00:00'
    const seconds = Math.floor(ms / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    const secs = seconds % 60
    
    if (hours > 0) {
      return `${String(hours).padStart(2, '0')}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
    }
    return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
  }

  const formatDateTime = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleString('ru-RU', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  const getMinBidAmount = () => {
    if (!auction) return 0
    if (!auction.current_bid) {
      return auction.start_price
    }
    return auction.current_bid.amount + auction.min_bid_step
  }

  const isBiddingEnabled = () => {
    return auction?.status === 'started' && auction?.status !== 'completed' && auction?.status !== 'cancelled'
  }

  if (loading) {
    return (
      <div className="auction-container">
        <div className="loading">Загрузка...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="auction-container">
        <div className="error">Ошибка: {error}</div>
        <button onClick={loadAuction}>Попробовать снова</button>
      </div>
    )
  }

  if (!auction) {
    return (
      <div className="auction-container">
        <div className="error">Аукцион не найден</div>
      </div>
    )
  }

  return (
    <div className="auction-container">
      <div className="auction-header">
        <div className="auction-title-section">
          <h1>Аукцион #{auction.id}</h1>
          <div className="user-info">
            <span className="user-label">Пользователь:</span>
            <span className="user-id">{userId}</span>
          </div>
        </div>
        <div className="auction-status">
          <span className="status-label">Статус:</span>
          <span className={`status status-${auction.status}`}>
            {auction.status === 'created' && 'Создан'}
            {auction.status === 'started' && 'Идет'}
            {auction.status === 'completed' && 'Завершен'}
            {auction.status === 'cancelled' && 'Отменен'}
          </span>
        </div>
      </div>

      <div className="auction-info">
        <div className="info-row">
          <span className="label">Лот:</span>
          <span className="value">{auction.lot_name || auction.lot_id}</span>
        </div>
        <div className="info-row">
          <span className="label">Начальная цена:</span>
          <span className="value">{auction.start_price} ₽</span>
        </div>
        <div className="info-row">
          <span className="label">Минимальный шаг:</span>
          <span className="value">{auction.min_bid_step} ₽</span>
        </div>
        {auction.start_time && (
          <div className="info-row">
            <span className="label">Время начала:</span>
            <span className="value">{formatDateTime(auction.start_time)}</span>
          </div>
        )}
        {auction.end_time && (
          <div className="info-row">
            <span className="label">Время завершения:</span>
            <span className="value">{formatDateTime(auction.end_time)}</span>
          </div>
        )}
        {auction.current_bid && (
          <div className="info-row">
            <span className="label">Текущая ставка:</span>
            <span className="value current-bid">{auction.current_bid.amount} ₽</span>
          </div>
        )}
        {timeRemaining !== null && timeRemaining > 0 && (
          <div className="info-row">
            <span className="label">Таймаут ставки:</span>
            <span className="value timeout">{formatTime(timeRemaining)}</span>
          </div>
        )}
        {!timeRemaining && auctionEndRemaining !== null && auctionEndRemaining > 0 && auction?.status === 'started' && (
          <div className="info-row">
            <span className="label">До завершения:</span>
            <span className="value auction-end">{formatTime(auctionEndRemaining)}</span>
          </div>
        )}
      </div>

      {auction.status === 'completed' && (
        <div className="auction-completed">
          <h2>Аукцион завершен!</h2>
          {auction.winner_id && (
            <div className="winner-info">
              <p>Победитель: {auction.winner_id}</p>
              <p className="final-price">Финальная цена: {auction.final_price} ₽</p>
            </div>
          )}
        </div>
      )}

      {auction.status === 'cancelled' && (
        <div className="auction-cancelled">
          <h2>Аукцион отменен</h2>
          <p>Ставки больше не принимаются</p>
          {cancellationReason && (
            <p className="cancellation-reason">
              <strong>Причина:</strong> {cancellationReason}
            </p>
          )}
        </div>
      )}

      {auction.status !== 'completed' && auction.status !== 'cancelled' && (
        <div className="bid-form-container">
          {bidError && (
            <div className="bid-error">
              <span className="error-icon">⚠️</span>
              <span className="error-message">{bidError}</span>
              <button 
                className="error-close" 
                onClick={() => setBidError(null)}
                aria-label="Закрыть"
              >
                ×
              </button>
            </div>
          )}
          <form onSubmit={handlePlaceBid} className="bid-form">
            <div className="form-group">
              <label>
                Минимальная ставка: {getMinBidAmount()} ₽
              </label>
              <input
                type="number"
                value={bidAmount}
                onChange={(e) => {
                  setBidAmount(e.target.value)
                  // Очищаем ошибку при изменении суммы
                  if (bidError) {
                    setBidError(null)
                  }
                }}
                placeholder={`Минимум ${getMinBidAmount()} ₽`}
                min={getMinBidAmount()}
                disabled={!isBiddingEnabled() || placingBid}
              />
            </div>
            <button
              type="submit"
              disabled={!isBiddingEnabled() || placingBid}
              className={!isBiddingEnabled() ? 'disabled' : ''}
            >
              {placingBid ? 'Отправка...' : isBiddingEnabled() ? 'Сделать ставку' : 'Ожидание старта'}
            </button>
          </form>
        </div>
      )}
    </div>
  )
}

export default AuctionPage

