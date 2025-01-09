document.addEventListener('DOMContentLoaded', function () {
  // Header dropdown
  document.addEventListener('click', function (event) {
    // Get the checkbox and other key elements
    const dropdownCheckbox = document.getElementById('header-user')
    const dropdownContainer = document.querySelector('.header-user')

    if (!dropdownCheckbox) return // Exit if the checkbox doesn't exist

    // Check if the click happened outside the dropdown
    const clickedInside = dropdownContainer.contains(event.target)

    if (!clickedInside) {
      // If the click is outside the dropdown, close it
      dropdownCheckbox.checked = false
    }
  })
  // Homepage countdown
  const countdownElement = document.getElementById('countdown')
  if (!countdownElement) return

  // Initialize time remaining from data attributes
  let hours = parseInt(countdownElement.getAttribute('data-hours'), 10)
  let minutes = parseInt(countdownElement.getAttribute('data-minutes'), 10)
  let seconds = parseInt(countdownElement.getAttribute('data-seconds'), 10)

  function updateCountdown() {
    // Decrement the time
    if (seconds > 0) {
      seconds--
    } else {
      if (minutes > 0) {
        minutes--
        seconds = 59
      } else {
        if (hours > 0) {
          hours--
          minutes = 59
          seconds = 59
        } else {
          // Countdown has finished
          clearInterval(countdownInterval)
          return
        }
      }
    }

    // Update the countdown display
    countdownElement.textContent =
      `${hours.toString().padStart(2, '0')}:` +
      `${minutes.toString().padStart(2, '0')}:` +
      `${seconds.toString().padStart(2, '0')}`
  }

  // Start the countdown
  updateCountdown() // Initial call to set the correct time immediately
  const countdownInterval = setInterval(updateCountdown, 1000)
})
