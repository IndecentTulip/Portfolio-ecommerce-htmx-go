(function() {
  let profdropdownButton = document.getElementById('ProfileDropdown');
  let profdropdownMenu = document.getElementById('ProfileDropdownMenu');

  let cartdropdownButton = document.getElementById('CartDropdown');
  let cartdropdownMenu = document.getElementById('CartDropdownMenu');

  profdropdownButton.addEventListener('click', () => {
    // Close the Cart dropdown if it's open
    if (!cartdropdownMenu.classList.contains('hidden')) {
      cartdropdownMenu.classList.add('hidden');
    }

    // Toggle the Profile dropdown menu
    profdropdownMenu.classList.toggle('hidden');
  });

  cartdropdownButton.addEventListener('click', () => {
    // Close the Profile dropdown if it's open
    if (!profdropdownMenu.classList.contains('hidden')) {
      profdropdownMenu.classList.add('hidden');
    }

    // Toggle the Cart dropdown menu
    cartdropdownMenu.classList.toggle('hidden');
  });

  window.addEventListener('click', (event) => {
    // Close both dropdowns if the click is outside
    if (!event.target.closest('.relative')) {
      profdropdownMenu.classList.add('hidden');
      cartdropdownMenu.classList.add('hidden');
    }
  });

})();
