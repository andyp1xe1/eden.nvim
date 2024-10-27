local M = {}

M.chan = nil

function M.start()
  if M.chan then
    return M.chan
  end
  M.chan = vim.fn.jobstart({ 'eden' }, {
    rpc = true,
    on_exit = function(_, code, _)
      if code ~= 0 then
        print("Eden Previewer exited with code: " .. code)
      end
      M.chan = nil
    end,
  })

  vim.fn.rpcnotify(M.chan, 'text_changed')
  return M.chan
end

function M.setup()
  vim.api.nvim_create_user_command('EdenStart', function()
    M.start()
    if not M.chan then
      print("Failed to start Eden.")
    end
  end, { nargs = 0 })

  vim.api.nvim_create_user_command('EdenEnd', function()
    if M.chan then
      vim.fn.jobstop(M.chan)
    end
  end, { nargs = 0 })

  vim.api.nvim_create_augroup('Eden', { clear = true })

  vim.api.nvim_create_autocmd({ 'TextChanged', 'TextChangedI', 'BufEnter' }, {
    group = 'Eden',
    pattern = '*',
    callback = function()
      if M.chan then
        vim.fn.rpcnotify(M.chan, 'text_changed')
      end
    end,
  })

  vim.api.nvim_create_autocmd({ 'CursorMoved', 'CursorMovedI' }, {
    group = 'Eden',
    pattern = '*',
    callback = function()
      if M.chan then
        vim.fn.rpcnotify(M.chan, 'scroll')
      end
    end,
  })
end

return M

--local function calculate_scroll_percentage()
--  local topline = vim.fn.line('w0')  -- First visible line
--  local total_lines = vim.api.nvim_buf_line_count(0)  -- Total number of lines in the buffer
--  local window_height = vim.api.nvim_win_get_height(0)  -- Current window height
--
--  -- Avoid division by zero
--  if total_lines == 0 or (total_lines - window_height) == 0 then
--    return 0
--  end
--
--  -- Calculate percentage scrolled
--  local percentage = ((topline - 1) / (total_lines - window_height)) * 100
--  return percentage
--end
