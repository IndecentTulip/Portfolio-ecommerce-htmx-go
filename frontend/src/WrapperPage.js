import Navbar from './Navbar';
import DropDowns from './DropDowns';
import { useState} from 'react';

const WrapperPage = ({ children }) => {

  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const [isCartOpen, setIsCartOpen] = useState(false);
  return (
  <div id="DivBody" className="bg-gray-100 text-gray-900 ">
    {/* change target for search, add a profile img */}
    <Navbar
      onProfileClick={()=>{
        setIsCartOpen(false);
        setIsProfileOpen((prev) => !prev);
        }
      }
      onCartClick={()=>{
        setIsProfileOpen(false);
        setIsCartOpen((prev) => !prev)
        }
      }
    />
    <DropDowns 
      isProfileOpen={isProfileOpen}
      isCartOpen={isCartOpen}
    />
    <div id="contentcontainer" class="contentcontainer p-6 mt-5">
      { children }
    </div>

  </div>
  );
};

export default WrapperPage;

