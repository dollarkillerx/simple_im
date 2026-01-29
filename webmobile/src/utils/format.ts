export function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const oneMinute = 60 * 1000
  const oneHour = 60 * oneMinute
  const oneDay = 24 * oneHour
  const oneWeek = 7 * oneDay

  if (diff < oneMinute) {
    return '刚刚'
  }

  if (diff < oneHour) {
    return `${Math.floor(diff / oneMinute)}分钟前`
  }

  if (diff < oneDay) {
    return `${Math.floor(diff / oneHour)}小时前`
  }

  if (diff < oneWeek) {
    return `${Math.floor(diff / oneDay)}天前`
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  if (year === now.getFullYear()) {
    return `${month}-${day}`
  }

  return `${year}-${month}-${day}`
}

export function formatChatTime(dateStr: string): string {
  const date = new Date(dateStr)
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes}`
}

export function formatFileSize(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`
  }

  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`
  }

  if (bytes < 1024 * 1024 * 1024) {
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

export function getDefaultAvatar(name: string): string {
  // Generate a data URL for a simple avatar with initials
  const initial = name.charAt(0).toUpperCase()
  const colors = [
    '#1989fa', '#07c160', '#ff976a', '#ee0a24',
    '#7232dd', '#1989fa', '#969799', '#ff6034',
  ]
  const colorIndex = name.charCodeAt(0) % colors.length
  const bgColor = colors[colorIndex]

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100">
      <rect width="100" height="100" fill="${bgColor}"/>
      <text x="50" y="50" font-family="Arial" font-size="40" fill="white"
            text-anchor="middle" dominant-baseline="central">${initial}</text>
    </svg>
  `

  return `data:image/svg+xml;base64,${btoa(svg)}`
}
