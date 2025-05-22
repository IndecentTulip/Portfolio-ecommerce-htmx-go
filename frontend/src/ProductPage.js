import { useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';

const ProductPage = () => {
  const productContainerRef = useRef(null);

  useEffect(() => {
    if (window.htmx && productContainerRef.current) {
      window.htmx.process(productContainerRef.current);
    }
  }, []);



  const { id } = useParams();
  return (
    <div ref={productContainerRef}>
      <div id="products" class="">
        <div id="test">test</div>
        <div><img  alt="Product Image"/></div>
        <h1>Turn this into API call </h1>
        <p>Prod Name: {id} </p>
        <div></div>
        <div></div>
        <div></div>
        <button 
          hx-put="/addtocart?id={{ .Values.Product.Id}}"
          hx-trigger="click"
          hx-swap="none"
        > Add to the cart</button>
          <div> you can add tags, ability to choose color if possible, etc</div>
      </div>
    </div>
  );
};

export default ProductPage;
