﻿using System;
using System.Windows.Forms;

namespace StreamDeckService
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
            string[] args = Environment.GetCommandLineArgs();
            foreach (String arg in args)
            {
                textBox1.Text += arg + "\r\n";
            }
        }

        private void cmQuit_Click(object sender, EventArgs e)
        {
            Application.Exit();
        }

        private void Form1_Activated(object sender, EventArgs e)
        {
            this.WindowState = FormWindowState.Minimized;
            if (this.WindowState == FormWindowState.Minimized)
            {
                Hide();
                notifyIcon1.Visible = true;
            }

        }
    }
}
