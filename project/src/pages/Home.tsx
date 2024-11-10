import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Search, Loader2 } from 'lucide-react';
import { Artist } from '../types';
import toast from 'react-hot-toast';

const Home = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    fetchArtists();
  }, []);

  const fetchArtists = async () => {
    try {
      const response = await fetch('http://localhost:8080/artists');
      const data = await response.json();
      setArtists(data);
    } catch (error) {
      toast.error('Failed to fetch artists');
    } finally {
      setLoading(false);
    }
  };

  const filteredArtists = artists.filter(artist =>
    artist.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-500" />
      </div>
    );
  }

  return (
    <div>
      <div className="max-w-2xl mx-auto mb-8">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search artists..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition-all"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredArtists.map((artist) => (
          <Link
            key={artist.id}
            to={`/artist/${artist.id}`}
            className="group bg-gray-800/50 rounded-lg overflow-hidden hover:transform hover:scale-105 transition-all duration-300 border border-gray-700/50 hover:border-indigo-500/50 hover:shadow-lg hover:shadow-indigo-500/10"
          >
            <div className="aspect-square overflow-hidden">
              <img
                src={artist.image}
                alt={artist.name}
                className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
              />
            </div>
            <div className="p-4">
              <h2 className="text-xl font-semibold mb-2">{artist.name}</h2>
              <p className="text-gray-400 text-sm">
                Created: {artist.creationDate}
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {artist.members.slice(0, 3).map((member, index) => (
                  <span
                    key={index}
                    className="text-xs px-2 py-1 bg-gray-700/50 rounded-full"
                  >
                    {member}
                  </span>
                ))}
                {artist.members.length > 3 && (
                  <span className="text-xs px-2 py-1 bg-gray-700/50 rounded-full">
                    +{artist.members.length - 3} more
                  </span>
                )}
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default Home;