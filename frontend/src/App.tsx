import React, { useState } from 'react';
import PlantsCRUD from './PlantsCRUD';
import PokemonCRUD from './PokemonCRUD';

// Define the UI interface
interface UI {
  name: string;
  content: React.ReactNode;
}

// Define the UIs object type
interface UIsType {
  [key: string]: UI;
}

const App: React.FC = () => {
  const [selectedUI, setSelectedUI] = useState<string>('rds');

  const UIs: UIsType = {
    rds: {
      name: 'Pokemon',
      content: <PokemonCRUD />
    },
    dynamo: {
      name: 'Plants',
      content: <PlantsCRUD />
    },
  };

  return (
    <div className="flex flex-col h-screen w-full">
      {/* Toggle Selector */}
      <div className="fixed top-0 left-0 right-0 z-50 bg-white shadow-md">
        <div className="relative bg-gray-200 rounded-full h-12 flex items-center max-w-md mx-auto my-2">
          <div
            className="absolute h-10 w-1/2 bg-white rounded-full shadow-md transition-all duration-300 ease-in-out"
            style={{
              left:
                selectedUI === 'rds' ? '0%' :
                  selectedUI === 'dynamo' ? '50%' :
                    '100%'
            }}
          />
          {Object.keys(UIs).map((ui) => (
            <button
              key={ui}
              className={`flex-1 flex items-center justify-center z-10 py-3 text-sm ${selectedUI === ui ? 'text-black font-bold' : 'text-gray-500'}`}
              onClick={() => setSelectedUI(ui)}
            >
              {UIs[ui].name}
            </button>
          ))}
        </div>
      </div>

      {/* Content Area - Full Screen with Top Padding */}
      <div className="flex-grow pt-16 h-full overflow-auto">
        {UIs[selectedUI].content}
      </div>
    </div>
  );
};

export default App;
