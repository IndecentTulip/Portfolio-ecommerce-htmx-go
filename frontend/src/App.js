import { SignedIn, SignedOut, RedirectToSignIn } from "@clerk/clerk-react";
import { BrowserRouter as Router, Routes, Route} from 'react-router-dom';
import Main from './Main';
import ProductPage from './ProductPage';
import WrapperPage from './WrapperPage';

const App = () => {
  return (
  <Router>
    <Routes>
      <Route path="/" element={
        <WrapperPage>
          <Main />
        </WrapperPage>
      } />
      <Route path="/p/:id" element={
        <WrapperPage>
          <ProductPage />
        </WrapperPage>
      } />
      <Route path="/c" element={
        <WrapperPage>
          <SignedIn>
            <div>.</div>
            <div>.</div>
            <div>cart page will be here</div>
          </SignedIn>
          <SignedOut>
            <RedirectToSignIn />
          </SignedOut>
        </WrapperPage>
      } />

    </Routes>
  </Router>
  );
};

export default App;

//  <body className='overflow-y-scroll'>

