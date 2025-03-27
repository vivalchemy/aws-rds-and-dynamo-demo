import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'; import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Trash2, Edit, RefreshCw } from 'lucide-react';
import { apiURL } from './lib/config';

const BACKEND_URL = apiURL + ':8080'

console.log(BACKEND_URL)

// Pokemon interface matching the backend
interface Pokemon {
  id?: number;
  name: string;
  type: string;
  hp: number;
  attack: number;
  defense: number;
  sp_attack: number;
  sp_defense: number;
  speed: number;
}

const PokemonCRUD: React.FC = () => {
  // State for form and data management
  const [pokemons, setPokemons] = useState<Pokemon[]>([]);
  const [currentPokemon, setCurrentPokemon] = useState<Pokemon>({
    name: '',
    type: '',
    hp: 0,
    attack: 0,
    defense: 0,
    sp_attack: 0,
    sp_defense: 0,
    speed: 0
  });
  const [isUpdateMode, setIsUpdateMode] = useState(false);

  // Fetch all Pokemon on component mount
  useEffect(() => {
    fetchPokemons();
  }, []);

  // Fetch Pokemon from backend
  const fetchPokemons = async () => {
    try {
      const response = await fetch(`${BACKEND_URL}/pokemon`);
      const data = await response.json();
      setPokemons(data);
    } catch (error) {
      console.error('Error fetching Pokemon:', error);
    }
  };

  // Handle input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCurrentPokemon(prev => ({
      ...prev,
      [name]: name.includes('hp') || name.includes('attack') ||
        name.includes('defense') || name.includes('speed')
        ? Number(value)
        : value
    }));
  };

  // Create or Update Pokemon
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const url = isUpdateMode ? `${BACKEND_URL}/pokemon/${currentPokemon.id}` : `${BACKEND_URL}/pokemon`;
    const method = isUpdateMode ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(currentPokemon)
      });

      if (response.ok) {
        fetchPokemons();
        // Reset form
        setCurrentPokemon({
          name: '',
          type: '',
          hp: 0,
          attack: 0,
          defense: 0,
          sp_attack: 0,
          sp_defense: 0,
          speed: 0
        });
        setIsUpdateMode(false);
      }
    } catch (error) {
      console.error('Error saving Pokemon:', error);
    }
  };

  // Prepare Pokemon for update
  const prepareUpdate = (pokemon: Pokemon) => {
    setCurrentPokemon(pokemon);
    setIsUpdateMode(true);
  };

  // Delete Pokemon
  const handleDelete = async (id: number) => {
    try {
      const response = await fetch(`${BACKEND_URL}/pokemon/${id}`, { method: 'DELETE' });
      if (response.ok) {
        fetchPokemons();
      }
    } catch (error) {
      console.error('Error deleting Pokemon:', error);
    }
  };

  return (
    <div className="flex p-4 space-x-4">
      {/* Left Side: Form */}
      <Card className="w-1/3">
        <CardHeader>
          <CardTitle>{isUpdateMode ? 'Update Pokemon' : 'Create Pokemon'}</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-2">
            <Input
              name="name"
              placeholder="Name"
              value={currentPokemon.name}
              onChange={handleInputChange}
              required
            />
            <Input
              name="type"
              placeholder="Type"
              value={currentPokemon.type}
              onChange={handleInputChange}
              required
            />
            <Input
              name="hp"
              type="number"
              placeholder="HP"
              value={currentPokemon.hp}
              onChange={handleInputChange}
              required
            />
            <Input
              name="attack"
              type="number"
              placeholder="Attack"
              value={currentPokemon.attack}
              onChange={handleInputChange}
              required
            />
            <Input
              name="defense"
              type="number"
              placeholder="Defense"
              value={currentPokemon.defense}
              onChange={handleInputChange}
              required
            />
            <Input
              name="sp_attack"
              type="number"
              placeholder="Special Attack"
              value={currentPokemon.sp_attack}
              onChange={handleInputChange}
              required
            />
            <Input
              name="sp_defense"
              type="number"
              placeholder="Special Defense"
              value={currentPokemon.sp_defense}
              onChange={handleInputChange}
              required
            />
            <Input
              name="speed"
              type="number"
              placeholder="Speed"
              value={currentPokemon.speed}
              onChange={handleInputChange}
              required
            />
            <Button type="submit" className="w-full">
              {isUpdateMode ? 'Update Pokemon' : 'Create Pokemon'}
            </Button>
            {isUpdateMode && (
              <Button
                type="button"
                variant="outline"
                className="w-full mt-2"
                onClick={() => {
                  setIsUpdateMode(false);
                  setCurrentPokemon({
                    name: '',
                    type: '',
                    hp: 0,
                    attack: 0,
                    defense: 0,
                    sp_attack: 0,
                    sp_defense: 0,
                    speed: 0
                  });
                }}
              >
                Cancel
              </Button>
            )}
          </form>
        </CardContent>
      </Card>

      {/* Right Side: Pokemon List */}
      <Card className="w-2/3">
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Pokemon List</CardTitle>
          <Button variant="outline" onClick={fetchPokemons}>
            <RefreshCw className="mr-2" /> Refresh
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>HP</TableHead>
                <TableHead>Attack</TableHead>
                <TableHead>Defense</TableHead>
                <TableHead>Sp. Atk</TableHead>
                <TableHead>Sp. Def</TableHead>
                <TableHead>Speed</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pokemons && pokemons.map(pokemon => (
                <TableRow key={pokemon.id}>
                  <TableCell>{pokemon.id}</TableCell>
                  <TableCell>{pokemon.name}</TableCell>
                  <TableCell>{pokemon.type}</TableCell>
                  <TableCell>{pokemon.hp}</TableCell>
                  <TableCell>{pokemon.attack}</TableCell>
                  <TableCell>{pokemon.defense}</TableCell>
                  <TableCell>{pokemon.sp_attack}</TableCell>
                  <TableCell>{pokemon.sp_defense}</TableCell>
                  <TableCell>{pokemon.speed}</TableCell>
                  <TableCell>
                    <div className="flex space-x-2">
                      <Button
                        size="icon"
                        variant="outline"
                        onClick={() => prepareUpdate(pokemon)}
                      >
                        <Edit size={16} />
                      </Button>
                      <Button
                        size="icon"
                        variant="destructive"
                        onClick={() => pokemon.id && handleDelete(pokemon.id)}
                      >
                        <Trash2 size={16} />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default PokemonCRUD;
