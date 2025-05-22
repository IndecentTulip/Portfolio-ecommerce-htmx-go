import { useEffect, useRef } from 'react';
import { useState } from 'react';

const Demo1 = () => {
  const [count, setCount] = useState(0);
  const [state, setState] = useState(true);

  const handleClick = () => {
    setCount(prevCount => prevCount + 1);
  };
  const handleState = () => {
    setState(prev => !prev);
  };

  const containerRef = useRef(null);

  useEffect(() => {
    if (window.htmx && containerRef.current) {
      window.htmx.process(containerRef.current);
    }
  }, []);
 // useEffect(() => {
 //   if (window.htmx && containerRef.current) {
 //     window.htmx.process(containerRef.current);
 //   }
 // }, [showHTMX]); //

  return (
    <div>
      <div className="max-w-xl mx-auto bg-white shadow-lg rounded-lg p-6">
        <button
          onClick={handleState}
          className="px-4 py-2 bg-red-600 hover:bg-blue-700 text-white rounded-md transition"
        >
          Update State and Brake HTMX 
        </button>


      </div>


    {state &&(
    <div ref={containerRef} className="space-y-4">
        <h3 className="text-2xl font-bold mb-4">React State Update</h3>
          <p className="mb-4 text-lg">
            State Updated: <strong>{count}</strong>
          </p>

        <button
          onClick={handleClick}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition"
        >
          Increment State
        </button>

      <h3 className="text-lg font-medium">HTMX</h3>
      <p id="content" className="text-gray-700 border p-2 rounded bg-gray-50">
        REPLACE ME
      </p>

      <button
        className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md"
        id="load-button"
        hx-get="/data/vars-1.html"
        hx-trigger="click"
        hx-target="#content"
      >
        Load New Content 
      </button>

      <button
        className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-md"
        id="add-btn"
        hx-get="/data/vars-1.html"
        hx-trigger="click"
        hx-target="#item-list"
      >
        Interact with another component
      </button>
    </div>
    )}
    </div>
  );
};

export default Demo1;

