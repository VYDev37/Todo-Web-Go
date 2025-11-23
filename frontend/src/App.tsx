import { useState, useEffect, useMemo } from "react";
import axios from "./config/axios";

import "./App.css";

interface Task { // Format: ngikut field backend
  ID: number;
  Name: string;
  Due: string;
  Done: boolean;

  // extra value (not for server-side)
  editing: boolean;
}

export default function App() {
  const [selection, setSelection] = useState<string>("show-all");
  const [loading, setLoading] = useState<boolean>(false);

  const [tasks, setTasks] = useState<Task[]>([]);
  const [body, setBody] = useState<Task>({
    ID: 0,
    Name: "",
    Due: "",
    Done: false,
    editing: false
  });

  const GetTasks = async () => {
    try {
      const response = await axios.get("todos");
      setTasks(response.data)
    } catch (err) {
      console.log(err);
      setTasks([]);
    }
  }

  const SetFormBody = (field: string, value: string | boolean | number) => {
    setBody(prev => ({ ...prev, [field]: value }));
  }

  const ModifyTask = (id: number, field: keyof Task, value: string | boolean) => {
    setTasks(prev =>
      prev.map(todo =>
        todo.ID === id ? { ...todo, [field]: value } : todo
      )
    );
  }

  const DeleteAllButton = async (e: React.FormEvent) => {
    e.preventDefault();
    setTasks([]);

    try {
      await axios.delete("todos");
    } catch (err) {
      console.log(err);
      return;
    }
  }

  const DeleteButton = async (e: React.FormEvent, id: number) => {
    e.preventDefault();
    if (loading)
      return;

    setLoading(true);
    try {
      await axios.delete(`todo/${id}`)
    } catch (err) {
      console.log(err);
      alert(err);
      return;
    } finally {
      setLoading(false);
    }
    GetTasks();
  }

  const EditButton = async (e: React.FormEvent, id: number) => {
    e.preventDefault();

    if (loading)
      return;

    const data = tasks.find(t => t.ID === id);
    if (!data)
      return;

    const isSaving = data.editing;
    if (isSaving) {
      setLoading(true);
      try {
        await axios.put(`todo/${data.ID}`, data);
      } catch (err) {
        console.log(err);
        alert(err);
        return;
      } finally {
        setLoading(false);
      }
    }

    setTasks(prev => prev.map(todo => todo.ID === id ? {...todo, editing: !todo.editing} : todo));
  };

  const SubmitButton = async (e: React.FormEvent) => {
    e.preventDefault();
    if (loading)
      return;

    if (!body.Name || !body.Due) {
      alert("Name and date must not be empty!");
      return;
    }

    const newid = tasks.length > 0 ? Math.max(...tasks.map(t => t.ID)) + 1 : 1;
    const newTask: Task = { ...body, ID: newid, editing: false };

    setTasks([...tasks, newTask]);
    setLoading(true);
    try {
      await axios.post("/add-todo", body);
      await GetTasks();
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }

    setBody({ ID: 0, Name: "", Due: "", Done: false, editing: false });
  }

  useEffect(() => {
    GetTasks();
  }, []);

  const filteredTasks = useMemo(() => {
    return tasks.filter(task => selection === "unfinished-only" ? !task.Done : true);
  }, [tasks, selection]);

  return (
    <div>
      <div className="flex justify-center items-center min-h-screen">
        <div className="bg-gray-100 w-full max-w-3xl rounded-xl shadow-lg p-8">
          <h1 className="font-bold text-4xl py-[1px] w-fit mx-auto text-gray-800 mb-8">Todo List</h1>

          <form className="flex flex-wrap items-center justify-center gap-4 mb-6" onSubmit={(e) => SubmitButton(e)}>
            <div className="relative min-w-[200px] flex-1">
              <input type="text" value={body.Name} onChange={(e) => SetFormBody("Name", e.target.value)} className="peer h-10 text-sm bg-gray-100 outline-none border border-gray-300 rounded px-3 py-1 w-full" placeholder=" " required />
              <label htmlFor="task-content" className="absolute left-3 top-2 text-sm text-gray-500 transition-all duration-300 
                            peer-placeholder-shown:top-2 peer-placeholder-shown:text-sm
                            peer-focus:top-[-0.5rem] peer-focus:text-xs peer-focus:text-blue-600 peer-[&amp;:not(:placeholder-shown)]:top-[-0.5rem] 
                            peer-[&amp;:not(:placeholder-shown)]:text-xs 
                            bg-gray-100 px-1">Task Name</label>
            </div>
            <div className="relative min-w-[160px] flex-1">
              <input type="date" value={body.Due} onChange={(e) => SetFormBody("Due", e.target.value)} className="peer h-10 text-sm bg-gray-100 text-gray-500 outline-none border border-gray-300 rounded px-3 py-1 w-full" placeholder="" required />
            </div>
            <div className="flex justify-center gap-4">
              <button disabled={loading} type="submit" className="px-3 py-1 rounded text-2xl bg-green-400 hover:bg-green-200 cursor-pointer transition-colors duration-300 ease-in-out">+</button>
            </div>
          </form>
          <div className="flex justify-between mb-4">
            <div className="flex gap-2">
              <select value={selection} onChange={(e) => setSelection(e.target.value)}
                className="border py-1 px-2 rounded">
                <option value="show-all">Show All</option>
                <option value="unfinished-only">Show Unfinished Only</option>
              </select>
              <button disabled={loading} className="btn-primary">FILTER</button>
            </div>
            <button disabled={loading} onClick={(e) => DeleteAllButton(e)} className="btn-danger">DELETE ALL</button>
          </div>
          <table className="table-auto w-full mt-5 border border-black border-collapse">
            <thead>
              <tr>
                <th className="table-cell">ID</th>
                <th className="table-cell">TASK</th>
                <th className="table-cell">DUE DATE</th>
                <th className="table-cell">STATUS</th>
                <th className="table-cell">ACTIONS</th>
              </tr>
            </thead>
            <tbody className="text-center">
              {
                filteredTasks.length > 0 ?
                  filteredTasks.map(task => (
                    <tr key={task.ID} className="hover:bg-gray-50">
                      <td className="table-cell">{task.ID}</td>
                      <td className="table-cell">
                        <input type="text" onChange={(e) => ModifyTask(task.ID, "Name", e.target.value)} value={task.Name} className="text-center" disabled={!task.editing} />
                      </td>
                      <td className="table-cell">
                        <input type="date" onChange={(e) => ModifyTask(task.ID, "Due", e.target.value)} value={task.Due} disabled={!task.editing} />
                      </td>
                      <td className="table-cell">
                        <input type="checkbox" onChange={(e) => ModifyTask(task.ID, "Done", e.target.checked)} checked={task.Done} disabled={!task.editing} />
                      </td>
                      <td className="table-cell">
                        <div className="mx-auto">
                          <button disabled={loading} onClick={(e) => EditButton(e, task.ID)}
                            className={`px-3 py-1 rounded text-white text-1xl font-bold transition-colors ${task.editing ? 'bg-green-500 hover:bg-green-600' : 'btn-primary'
                              }`}>{task.editing ? "SAVE" : "EDIT"}</button>
                          <button disabled={loading} onClick={(e) => DeleteButton(e, task.ID)} className="btn-danger mx-1">DELETE</button>
                        </div>
                      </td>
                    </tr>
                  )) : (
                    <tr>
                      <td colSpan={5} className="px-2 py-1">No data available.</td>
                    </tr>
                  )
              }
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}