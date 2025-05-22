import './Main.css';
import { useEffect, useRef } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';

const Main = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const productContainerRef = useRef(null);

  // Read search params first
  const search = searchParams.get('search') ?? '';
  const num = searchParams.get('num') ?? '0'; // fallback to 0

  // Build mainLink based on search params
  const mainLink = search === ''
    ? `http://localhost:25000/main?start=0&num=${num}`
    : `http://localhost:25000/main?start=0&search=${encodeURIComponent(search)}&num=${num}`;

  useEffect(() => {
    // Ensure ?num=0 is in the URL
    if (!searchParams.has('num')) {
      const newParams = new URLSearchParams(searchParams);
      newParams.set('num', '0');
      navigate(`/?${newParams.toString()}`, { replace: true });
      return; // wait for navigation to complete
    }

    // Inject HTMX component and process
    if (window.htmx && productContainerRef.current) {
      productContainerRef.current.innerHTML = `
        <div 
          hx-get="${mainLink}"
          hx-trigger="revealed"
          hx-swap="outerHTML" class="text-center p-4 bg-gray-200 rounded-lg shadow-md hover:bg-gray-300 transition"
        ></div>
      `;
      window.htmx.process(productContainerRef.current);
    }
  }, [search, num, navigate, searchParams, mainLink]);

  return (
    <div>
      <div
        id="products"
        ref={productContainerRef}
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 p-4"
      >
        {/* HTMX div will be injected here */}
      </div>


      <div id="nextpage" classname="col-auto row-auto">
        <div class="flex justify-center space-x-2 mt-4">
          <button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition">
            &#8203;
          </button>
          <button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition">
            &#8203;
          </button>
          <button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition">
            &#8203;
          </button>
          <button class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition">
            &#8203;
          </button>
        </div>
      </div>

    </div>
  );
};

export default Main;

