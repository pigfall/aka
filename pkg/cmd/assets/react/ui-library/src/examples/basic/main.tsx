import React from 'react';
import { createRoot } from 'react-dom/client';
import  Button  from '../../components/Button';

const container = document.getElementById('root')!
const root = createRoot(container);
root.render(<Button>test</Button>);
