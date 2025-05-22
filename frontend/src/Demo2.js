import { useEffect, useRef } from 'react';

const Demo2 = () => {

  const containerRef = useRef(null);
  useEffect(() => {
    if (window.htmx && containerRef.current) {
      window.htmx.process(containerRef.current);
    }
  }, []);

  return (
    <div ref={containerRef}>
      <h3>Another Component:</h3>

      <ul id="item-list">
      </ul>
    </div>
  );
}


export default Demo2;
