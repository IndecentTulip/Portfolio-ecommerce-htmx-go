//import { useEffect, useRef } from 'react';
import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

const Navbar = ({onProfileClick, onCartClick}) => {
  //const containerRef = useRef(null);
  const [searchTerm, setSearchTerm] = useState('');
  //const [searchParams] = useSearchParams();

  //if (searchTerm == ''){
  //  setSearchTerm( searchParams.get('search') ?? '');
  //}
  
  const navigate = useNavigate();
  //const location = useLocation();

  //useEffect(() => {
  //  if (window.htmx && containerRef.current) {
  //    window.htmx.process(containerRef.current);
  //  }
  //}, []);

  const handleImageError = (e) => {
    e.target.onerror = null; // Prevents infinite loop
    e.target.src = process.env.PUBLIC_URL + '/default_profile.jpg';
  };
  const src = "google.com"

  const cartImg = process.env.PUBLIC_URL + '/Cart.png'

  //const handleKeyDown = (e) => {
  //  if (e.key === 'Enter') {
  //    const params = new URLSearchParams(location.search);
  //    params.set('search', searchTerm); // Add or replace the 'search' parameter
  //
  //    // Keep existing path and add updated query string
  //    navigate(`${location.pathname}?${params.toString()}`);
  //  }
  //};
  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
       // Navigate to "/" and include the search term as a query parameter
       navigate(`/?num=0&search=${encodeURIComponent(searchTerm)}`);
    }
  }


  // ref={containerRef}
  return (
    <div className="fixed top-0 left-0 right-0 bg-white shadow-md z-50 grid grid-cols-3 items-center px-4">

      {/*>burger*/}
      <div>
        put burger menu
        <a href='/'>
        put logo here
        </a>
      </div>
      {/*<burger*/}

      {/*>search bar*/}
      <input
        classNameName="w-80"
        type="search"
        name="search"
        placeholder="Type To Search"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        onKeyDown={handleKeyDown}
      />
      {/*<search bar*/}

      {/*>clickable elements in a navbar  */}
      <div className="flex justify-end gap-10 px-2">

        <div className="relative">
          <button 
            className="bg-blue-500 rounded-md"
            id="ProfileDropdown"
            onClick={() => onProfileClick()}
          >
            <img 
              type="button" 
              src={src} alt="profile image" 
              className="w-12 h-12 rounded-full object-cover"
              onError={handleImageError}
            />
          </button>
        </div>
        <div className="relative">
          <button 
            className="bg-blue-500 rounded-md" 
            id="CartDropdown"
            onClick={() => onCartClick()}
          >
            <img className="w-14" src={cartImg} />
          </button>
        </div>

      </div>
      {/*<clickable elements in a navbar  */}

    </div>
  );
};

export default Navbar;

