import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Trash2, Edit, RefreshCw } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { apiURL } from './lib/config';

const BACKEND_URL = apiURL + ':8081'

console.log(BACKEND_URL)


// Plant interface matching the backend
interface Plant {
  id?: string;
  name: string;
  scientific_name: string;
  family: string;
  type: string;
  sunlight_required: string;
  water_interval: number;
  height: number;
  native: string;
  indoor: boolean;
}


const PlantsCRUD: React.FC = () => {
  // State for form and data management
  const [plants, setPlants] = useState<Plant[]>([]);
  const [currentPlant, setCurrentPlant] = useState<Plant>({
    name: '',
    scientific_name: '',
    family: '',
    type: '',
    sunlight_required: '',
    water_interval: 0,
    height: 0,
    native: '',
    indoor: false
  });
  const [isUpdateMode, setIsUpdateMode] = useState(false);

  // Fetch all Plants on component mount
  useEffect(() => {
    fetchPlants();
  }, []);

  // Fetch Plants from backend
  const fetchPlants = async () => {
    try {
      const response = await fetch(`${BACKEND_URL}/plants`);
      const data = await response.json();
      setPlants(data);
    } catch (error) {
      console.error('Error fetching Plants:', error);
    }
  };

  // Handle input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setCurrentPlant(prev => ({
      ...prev,
      [name]: type === 'checkbox'
        ? checked
        : (name === 'water_interval' || name === 'height')
          ? Number(value)
          : value
    }));
  };

  // Create or Update Plant
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const url = isUpdateMode ? `${BACKEND_URL}/plants/${currentPlant.id}` : `${BACKEND_URL}/plants`;
    const method = isUpdateMode ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(currentPlant)
      });

      if (response.ok) {
        fetchPlants();
        // Reset form
        setCurrentPlant({
          name: '',
          scientific_name: '',
          family: '',
          type: '',
          sunlight_required: '',
          water_interval: 0,
          height: 0,
          native: '',
          indoor: false
        });
        setIsUpdateMode(false);
      }
    } catch (error) {
      console.error('Error saving Plant:', error);
    }
  };

  // Prepare Plant for update
  const prepareUpdate = (plant: Plant) => {
    setCurrentPlant(plant);
    setIsUpdateMode(true);
  };

  // Delete Plant
  const handleDelete = async (id: string) => {
    try {
      const response = await fetch(`${BACKEND_URL}/plants/${id}`, { method: 'DELETE' });
      if (response.ok) {
        fetchPlants();
      }
    } catch (error) {
      console.error('Error deleting Plant:', error);
    }
  };

  return (
    <div className="flex p-4 space-x-4 bg-slate-50 min-h-screen">
      {/* Left Side: Form */}
      <Card className="w-1/3">
        <CardHeader>
          <CardTitle>{isUpdateMode ? 'Update Plant' : 'Create Plant'}</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-2">
            <Input
              name="name"
              placeholder="Name"
              value={currentPlant.name}
              onChange={handleInputChange}
              required
            />
            <Input
              name="scientific_name"
              placeholder="Scientific Name"
              value={currentPlant.scientific_name}
              onChange={handleInputChange}
              required
            />
            <Input
              name="family"
              placeholder="Family"
              value={currentPlant.family}
              onChange={handleInputChange}
              required
            />
            <Input
              name="type"
              placeholder="Type"
              value={currentPlant.type}
              onChange={handleInputChange}
              required
            />
            <Input
              name="sunlight_required"
              placeholder="Sunlight Required"
              value={currentPlant.sunlight_required}
              onChange={handleInputChange}
              required
            />
            <Input
              name="water_interval"
              type="number"
              placeholder="Water Interval (days)"
              value={currentPlant.water_interval}
              onChange={handleInputChange}
              required
            />
            <Input
              name="height"
              type="number"
              step="0.1"
              placeholder="Height (meters)"
              value={currentPlant.height}
              onChange={handleInputChange}
              required
            />
            <Input
              name="native"
              placeholder="Native Region"
              value={currentPlant.native}
              onChange={handleInputChange}
              required
            />
            <div className="flex items-center space-x-2">
              <Switch
                id="indoor"
                name="indoor"
                checked={currentPlant.indoor}
                onCheckedChange={(checked) => setCurrentPlant(prev => ({
                  ...prev,
                  indoor: checked
                }))}
              />
              <Label htmlFor="indoor">Indoor Plant</Label>
            </div>
            <Button type="submit" className="w-full mt-2">
              {isUpdateMode ? 'Update Plant' : 'Create Plant'}
            </Button>
            {isUpdateMode && (
              <Button
                type="button"
                variant="outline"
                className="w-full mt-2"
                onClick={() => {
                  setIsUpdateMode(false);
                  setCurrentPlant({
                    name: '',
                    scientific_name: '',
                    family: '',
                    type: '',
                    sunlight_required: '',
                    water_interval: 0,
                    height: 0,
                    native: '',
                    indoor: false
                  });
                }}
              >
                Cancel
              </Button>
            )}
          </form>
        </CardContent>
      </Card>

      {/* Right Side: Plants List */}
      <Card className="w-2/3">
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Plants List</CardTitle>
          <Button variant="outline" onClick={fetchPlants}>
            <RefreshCw className="mr-2" /> Refresh
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Scientific Name</TableHead>
                <TableHead>Family</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Sunlight</TableHead>
                <TableHead>Water Interval</TableHead>
                <TableHead>Height</TableHead>
                <TableHead>Native</TableHead>
                <TableHead>Indoor</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {plants && plants.map(plant => (
                <TableRow key={plant.id}>
                  <TableCell>{plant.name}</TableCell>
                  <TableCell>{plant.scientific_name}</TableCell>
                  <TableCell>{plant.family}</TableCell>
                  <TableCell>{plant.type}</TableCell>
                  <TableCell>{plant.sunlight_required}</TableCell>
                  <TableCell>{plant.water_interval} days</TableCell>
                  <TableCell>{plant.height} m</TableCell>
                  <TableCell>{plant.native}</TableCell>
                  <TableCell>{plant.indoor ? 'Yes' : 'No'}</TableCell>
                  <TableCell>
                    <div className="flex space-x-2">
                      <Button
                        size="icon"
                        variant="outline"
                        onClick={() => prepareUpdate(plant)}
                      >
                        <Edit size={16} />
                      </Button>
                      <Button
                        size="icon"
                        variant="destructive"
                        onClick={() => plant.id && handleDelete(plant.id)}
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

export default PlantsCRUD;
