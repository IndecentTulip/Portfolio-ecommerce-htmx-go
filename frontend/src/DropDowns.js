import CartItems from './CartItems';

const DropDowns = ({ isProfileOpen, isCartOpen }) => {
  return (
    <div className="fixed top-16 right-0 flex justify-end w-full gap-6 z-50">
      {/*>Profile Dropdown */}
      {isProfileOpen && (
        <div className="absolute bg-gray-200 font-semibold text-lg px-4 py-2 mt-2 mr-11 rounded shadow-md">
          <span className="block">Welcome:</span>
          <hr />
          <span className="block mt-2">Login With</span>
          <a
            href="http://portfolio.serverpit.com:25000/login?type=google"
            className="block bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 mt-2 transition"
          >
            Google
          </a>
          <a
            href="http://portfolio.serverpit.com:25000/login?type=github"
            className="block bg-gray-800 text-white py-2 px-4 rounded-lg hover:bg-gray-900 mt-2 transition"
          >
            GitHub
          </a>
        </div>
      )}
      {/*<Profile Dropdown */}

      {/*>Cart Dropdown */}
      {isCartOpen && (
        <div className="absolute bg-gray-200 font-semibold text-lg px-4 py-2 mt-2 overflow-y-scroll max-h-64 mr-5 rounded shadow-md">
          <div>Items</div>
          <button
            className="bg-green-500 rounded-md mt-2 px-4 py-2 text-white hover:bg-green-600"
            onClick={() => (window.location.href = '/c/')}
          >
            Checkout
          </button>
          <div id="items" className="space-y-4 mt-4">
            <CartItems />
          </div>
        </div>
      )}
      {/*<Cart Dropdown */}
    </div>
  );
};

export default DropDowns;

